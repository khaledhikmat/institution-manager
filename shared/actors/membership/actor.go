package membershipactor

import (
	"context"
	"fmt"

	"github.com/dapr/go-sdk/actor"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/khaledhikmat/campaign-manager/shared/service/member"
	"github.com/khaledhikmat/campaign-manager/shared/service/membership"
)

// Injected DAPR client and other services
var Daprclient dapr.Client
var MembershipService membership.IService

func MembershipActorFactory() actor.ServerContext {
	return &MembershipActor{}
}

/*
Setup a timer to update the domain database if active
Setup a reminder to update the exchange rate if active
Setup a timer to push or do a realtime pusher
*/

type MembershipActor struct {
	// Platform dependencies
	actor.ServerImplBaseCtx

	// Misc State
	miscState miscState

	// Main State
	mainState membership.Membership

	// Transactions State
	transactionsState []member.MemberTransaction
}

func (a *MembershipActor) Type() string {
	return "MembershipActorType"
}

func (a *MembershipActor) Main(ctx context.Context) (membership.Membership, error) {
	fmt.Println("MembershipActor", "Main")
	a.getState(ctx, mainStateKey)
	return a.mainState, nil
}

func (a *MembershipActor) Transactions(ctx context.Context) ([]member.MemberTransaction, error) {
	fmt.Println("MembershipActor", "Transactions")
	a.getState(ctx, transactionsStateKey)
	return a.transactionsState, nil
}

func (a *MembershipActor) Update(ctx context.Context, c membership.Membership) error {
	fmt.Println("MembershipActor", "Update")
	a.getState(ctx, mainStateKey)

	// Update the actor state with the member update-able data element
	a.mainState.InstitutionID = c.InstitutionID
	a.mainState.Name = c.Name
	a.mainState.Description = c.Description
	a.mainState.Benefits = c.Benefits
	a.saveState(ctx, mainStateKey)

	// Update the database
	err := MembershipService.UpdateMembership(&c)
	if err != nil {
		return err
	}

	return nil
}

func (a *MembershipActor) Transact(ctx context.Context, evt member.MemberTransaction) error {
	fmt.Printf("MembershipActor Transact - MEMBERSHIP ID: %s - MEMBER ID: %s - AMOUNT: %d\n", evt.MembershipID, evt.MemberID, evt.Amount)
	a.getState(ctx, miscStateKey)
	a.getState(ctx, mainStateKey)
	evt.ExgRate = a.miscState.ExchangeRate

	a.mainState.Transactions++
	if evt.Type == "Purchase" {
		a.mainState.Purchases++
	} else if evt.Type == "Renewal" {
		a.mainState.Renewals++
	} else if evt.Type == "Canx" {
		a.mainState.Cancellations++
	}
	a.saveState(ctx, mainStateKey)

	a.getState(ctx, transactionsStateKey)
	// TODO: The actor state should limit the transactions to the last 10 for example
	a.transactionsState = append(a.transactionsState, evt)
	a.saveState(ctx, transactionsStateKey)

	return nil
}
