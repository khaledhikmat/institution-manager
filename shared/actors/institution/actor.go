package institutionactor

import (
	"context"
	"fmt"
	"time"

	"github.com/dapr/go-sdk/actor"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/khaledhikmat/institution-manager/shared/service/institution"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
)

// Injected DAPR client and other services
var Daprclient dapr.Client
var InstitutionService institution.IService

func InstitutionActorFactory() actor.ServerContext {
	return &InstitutionActor{}
}

/*
Setup a timer to update the domain database if active
Setup a reminder to update the exchange rate if active
*/

type InstitutionActor struct {
	// Platform dependencies
	actor.ServerImplBaseCtx

	// Main State
	mainState institution.Institution

	// Misc State
	miscState miscState
}

func (a *InstitutionActor) Type() string {
	return "InstitutionActorType"
}

func (a *InstitutionActor) Main(ctx context.Context) (institution.Institution, error) {
	fmt.Println("InstitutionActor", "Main")
	a.getState(ctx, mainStateKey)
	return a.mainState, nil
}

func (a *InstitutionActor) Update(ctx context.Context, c institution.Institution) error {
	fmt.Println("InstitutionActor", "Update")
	a.getState(ctx, mainStateKey)

	// Update the actor state with the institution update-able data element
	a.mainState.Name = c.Name
	a.mainState.City = c.City
	a.mainState.State = c.State
	a.mainState.Country = c.Country
	a.saveState(ctx, mainStateKey)

	// Update the database
	err := InstitutionService.UpdateInstitution(&c)
	if err != nil {
		return err
	}

	return nil
}

func (a *InstitutionActor) Pledge(ctx context.Context, evt member.MemberPledge) error {
	fmt.Printf("InstitutionActor pledge - CAMPAIGN ID: %s - MEMBER ID: %s - AMOUNT: %d\n", evt.CampaignID, evt.MemberID, evt.Amount)
	a.getState(ctx, miscStateKey)
	a.getState(ctx, mainStateKey)
	evt.Time = time.Now()
	evt.ExgRate = a.miscState.ExchangeRate

	a.mainState.Pledges++
	a.mainState.Donors++ // TODO: Not accurate...must be unique
	a.saveState(ctx, mainStateKey)

	return nil
}
