package memberactor

import (
	"context"
	"fmt"
	"time"

	"github.com/dapr/go-sdk/actor"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/khaledhikmat/campaign-manager/shared/equates"
	"github.com/khaledhikmat/campaign-manager/shared/service/institution"
	"github.com/khaledhikmat/campaign-manager/shared/service/member"
)

// Injected DAPR client and other services
var Daprclient dapr.Client
var MemberService member.IService

func MemberActorFactory() actor.ServerContext {
	return &MemberActor{}
}

/*
Setup a timer to update the domain database if active
Setup a reminder to update the exchange rate if active
Setup a timer to push or do a realtime pusher
*/

type MemberActor struct {
	// Platform dependencies
	actor.ServerImplBaseCtx

	// Misc State
	miscState miscState

	// Main State
	mainState member.Member

	// Instititions State
	institutionsState []institution.Institution

	// Memberships State
	membershipsState []member.MemberMembership

	// Transactions State
	transactionsState []member.MemberTransaction

	// Pledges State
	pledgesState []member.MemberPledge
}

func (a *MemberActor) Type() string {
	return "MemberActorType"
}

func (a *MemberActor) Main(ctx context.Context) (member.Member, error) {
	fmt.Println("MemberActor", "Main")
	a.getState(ctx, mainStateKey)
	return a.mainState, nil
}

func (a *MemberActor) Institutions(ctx context.Context) ([]institution.Institution, error) {
	fmt.Println("MemberActor", "Institutions")
	a.getState(ctx, institutionsStateKey)
	return a.institutionsState, nil
}

func (a *MemberActor) Memberships(ctx context.Context) ([]member.MemberMembership, error) {
	fmt.Println("MemberActor", "Memberships")
	a.getState(ctx, membershipsStateKey)
	return a.membershipsState, nil
}

func (a *MemberActor) Transactions(ctx context.Context) ([]member.MemberTransaction, error) {
	fmt.Println("MemberActor", "Transactions")
	a.getState(ctx, transactionsStateKey)
	return a.transactionsState, nil
}

func (a *MemberActor) Pledges(ctx context.Context) ([]member.MemberPledge, error) {
	fmt.Println("MemberActor", "Pledges")
	a.getState(ctx, pledgesStateKey)
	return a.pledgesState, nil
}

func (a *MemberActor) Update(ctx context.Context, c member.Member) error {
	fmt.Println("MemberActor", "Update")
	a.getState(ctx, mainStateKey)

	// Update the actor state with the member update-able data element
	a.mainState.Role = c.Role
	a.mainState.Type = c.Type
	a.mainState.Parent = c.Parent
	a.mainState.ExternalID = c.ExternalID
	a.mainState.Name = c.Name
	a.mainState.Email = c.Email
	a.mainState.Phone = c.Phone
	a.saveState(ctx, mainStateKey)

	// Update the database
	err := MemberService.UpdateMember(&c)
	if err != nil {
		return err
	}

	return nil
}

func (a *MemberActor) Transact(ctx context.Context, evt member.MemberTransaction) error {
	fmt.Printf("MemberActor Transact - MEMBERSHIP ID: %s - MEMBER ID: %s - AMOUNT: %d\n", evt.MembershipID, evt.MemberID, evt.Amount)
	a.getState(ctx, miscStateKey)
	a.getState(ctx, mainStateKey)
	evt.ExgRate = a.miscState.ExchangeRate

	// TODO: The actor will process transaction against the member.
	// If no error, it will proceed to affect the state and publish the transaction event
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

	// If there is no error, publish the event
	err := Daprclient.PublishEvent(ctx, equates.CAMPAIGN_PUB_SUB, equates.TRANSACTIONS_TOPIC, evt)
	if err != nil {
		fmt.Printf("publish event to pledges topic errored out %v\n", err)
		return err
	}

	return nil
}

func (a *MemberActor) Pledge(ctx context.Context, evt member.MemberPledge) error {
	fmt.Printf("MemberActor pledge - CAMPAIGN ID: %s - MEMBER ID: %s - AMOUNT: %d\n", evt.CampaignID, evt.MemberID, evt.Amount)
	a.getState(ctx, miscStateKey)
	a.getState(ctx, mainStateKey)
	evt.Time = time.Now()
	evt.ExgRate = a.miscState.ExchangeRate

	a.mainState.Pledges++
	a.saveState(ctx, mainStateKey)

	a.getState(ctx, pledgesStateKey)
	// TODO: The actor state should limit the pledges to the last 10 for example
	a.pledgesState = append(a.pledgesState, evt)
	a.saveState(ctx, pledgesStateKey)

	return nil
}
