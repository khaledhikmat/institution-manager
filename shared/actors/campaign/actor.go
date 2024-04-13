package campaignactor

import (
	"context"
	"fmt"
	"time"

	"github.com/dapr/go-sdk/actor"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/khaledhikmat/institution-manager/shared/equates"
	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
)

// Injected DAPR client and other services
var Daprclient dapr.Client
var CampaignService campaign.IService
var MemberService member.IService

func CampaignActorFactory() actor.ServerContext {
	return &CampaignActor{}
}

/*
Setup a timer to update the domain database if active
Setup a reminder to update the exchange rate if active
Setup a timer to push or do a realtime pusher
*/

type CampaignActor struct {
	// Platform dependencies
	actor.ServerImplBaseCtx

	// Misc State
	miscState miscState

	// Main State
	mainState campaign.Campaign

	// Pledges State
	pledgesState []member.MemberPledge
}

func (a *CampaignActor) Type() string {
	return "CampaignActorType"
}

func (a *CampaignActor) Main(ctx context.Context) (campaign.Campaign, error) {
	fmt.Println("CampaignActor", "Main")
	a.getState(ctx, mainStateKey)
	return a.mainState, nil
}

func (a *CampaignActor) Pledges(ctx context.Context) ([]member.MemberPledge, error) {
	fmt.Println("CampaignActor", "Pledges")
	a.getState(ctx, pledgesStateKey)
	return a.pledgesState, nil
}

func (a *CampaignActor) Update(ctx context.Context, c campaign.Campaign) error {
	fmt.Println("CampaignActor", "Update")
	a.getState(ctx, mainStateKey)

	// Update the actor state with the campaign update-able data element
	a.mainState.Name = c.Name
	a.mainState.Description = c.Description
	a.mainState.Virtual = c.Virtual
	a.mainState.ImageURL = c.ImageURL
	a.mainState.Goal = c.Goal
	a.mainState.Currency = c.Currency
	a.mainState.Duration = c.Duration
	a.saveState(ctx, mainStateKey)

	// Update the database
	err := CampaignService.UpdateCampaign(&c)
	if err != nil {
		return err
	}

	return nil
}

func (a *CampaignActor) Pledge(ctx context.Context, evt member.MemberPledge) error {
	fmt.Printf("CampaignActor pledge - CAMPAIGN ID: %s - MEMBER ID: %s - AMOUNT: %d\n", evt.CampaignID, evt.MemberID, evt.Amount)
	a.getState(ctx, miscStateKey)
	a.getState(ctx, mainStateKey)
	evt.Time = time.Now()
	evt.ExgRate = a.miscState.ExchangeRate

	// TODO: The actor will process pledge against the campaign.
	// If no error, it will proceed to affect the state and publish the pledge event
	a.mainState.Pledges++
	a.mainState.TotPledgeAmount += evt.Amount
	a.mainState.AwayFromGoal = a.mainState.Goal - a.mainState.TotPledgeAmount
	if a.mainState.MaxPledgeAmount <= evt.Amount {
		a.mainState.MaxPledgeAmount = evt.Amount
	}
	if a.mainState.MinPledgeAmount >= evt.Amount {
		a.mainState.MinPledgeAmount = evt.Amount
	}
	if a.mainState.Pledges > 0 {
		a.mainState.AvgPledgeAmount = a.mainState.TotPledgeAmount / a.mainState.Pledges
	}
	a.mainState.TotPledgeAmountUSD += evt.Amount * int64(evt.ExgRate)
	a.mainState.Donors++ // TODO: Not accurate...it must be unique donors
	//MostGenerousDonor  user.User `json:"user"`
	a.saveState(ctx, mainStateKey)

	a.getState(ctx, pledgesStateKey)
	// TODO: The actor state should limit the pledges to the last 10 for example
	a.pledgesState = append(a.pledgesState, evt)
	a.saveState(ctx, pledgesStateKey)

	// If there is no error in processing the campaign,
	// the actor fires two events: a pledge event and a campaign event
	err := Daprclient.PublishEvent(ctx, equates.CAMPAIGN_PUB_SUB, equates.PLEDGES_TOPIC, evt)
	if err != nil {
		fmt.Printf("publish event to pledges topic errored out %v\n", err)
		return err
	}

	// The campaign event is for the externalizer
	evt2 := equates.CampaignPledges{
		Campaign: a.mainState,
		Pledges:  a.pledgesState,
	}
	err = Daprclient.PublishEvent(ctx, equates.CAMPAIGN_PUB_SUB, equates.CAMPAIGNS_TOPIC, evt2)
	if err != nil {
		fmt.Printf("publish event to campaigns topic errored out %v\n", err)
		return err
	}

	return nil
}

func (a *CampaignActor) Confirm(ctx context.Context, evt member.MemberPledge) error {
	fmt.Printf("CampaignActor Confirm - CAMPAIGN ID: %s - MEMBER ID: %s - AMOUNT: %d\n", evt.CampaignID, evt.MemberID, evt.Amount)
	// TODO:
	return nil
}

func (a *CampaignActor) Decline(ctx context.Context, evt member.MemberPledge) error {
	fmt.Printf("CampaignActor Decline - CAMPAIGN ID: %s - MEMBER ID: %s - AMOUNT: %d\n", evt.CampaignID, evt.MemberID, evt.Amount)
	// TODO:
	return nil
}
