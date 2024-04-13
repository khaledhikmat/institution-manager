package flow

import (
	"fmt"

	"github.com/dapr/go-sdk/workflow"
	campaignactor "github.com/khaledhikmat/institution-manager/shared/actors/campaign"
)

func submitPledgeActivity(ctx workflow.ActivityContext) (any, error) {
	var input CampaignPledge
	if err := ctx.GetInput(&input); err != nil {
		fmt.Printf("submitPledgeActivity - error: %v\n", err)
		return "", err
	}

	fmt.Printf("submitPledgeActivity - input: %v\n", input)

	// Resolve actor by campaign id
	campaignActorProxy := campaignactor.NewCampaignActor(input.Campaign.ID)
	DaprClient.ImplActorClientStub(campaignActorProxy)

	// Call campaign actor to submit pledge
	err := campaignActorProxy.Pledge(ctx.Context(), input.Pledge)
	if err != nil {
		fmt.Printf("submitPledgeActivity - error2: %v\n", err)
		return input, err
	}

	fmt.Printf("submitPledgeActivity - output: %v\n", input)
	return input, nil
}

func confirmPledgeActivity(ctx workflow.ActivityContext) (any, error) {
	var input CampaignPledge
	if err := ctx.GetInput(&input); err != nil {
		fmt.Printf("confirmPledgeActivity - error: %v\n", err)
		return "", err
	}

	fmt.Printf("confirmPledgeActivity - input: %v\n", input)

	// Resolve actor by campaign id
	campaignActorProxy := campaignactor.NewCampaignActor(input.Campaign.ID)
	DaprClient.ImplActorClientStub(campaignActorProxy)

	// Call campaign actor to confirm pledge
	err := campaignActorProxy.Confirm(ctx.Context(), input.Pledge)
	if err != nil {
		fmt.Printf("confirmPledgeActivity - error2: %v\n", err)
		return input, err
	}

	fmt.Printf("confirmPledgeActivity - output: %v\n", input)
	return input, nil
}

func declinePledgeActivity(ctx workflow.ActivityContext) (any, error) {
	var input CampaignPledge
	if err := ctx.GetInput(&input); err != nil {
		fmt.Printf("declinePledgeActivity - error: %v\n", err)
		return "", err
	}

	fmt.Printf("declinePledgeActivity - input: %v\n", input)

	// Resolve actor by campaign id
	campaignActorProxy := campaignactor.NewCampaignActor(input.Campaign.ID)
	DaprClient.ImplActorClientStub(campaignActorProxy)

	// Call campaign actor to decline pledge
	err := campaignActorProxy.Decline(ctx.Context(), input.Pledge)
	if err != nil {
		fmt.Printf("declinePledgeActivity - error2: %v\n", err)
		return input, err
	}

	// Resolve actor by pledge id
	// pledgeActorProxy := pledgeactor.NewPledgeActor(input.Pledge.ID)
	// DaprClient.ImplActorClientStub(pledgeActorProxy)

	// // Call pledge actor to decline pledge
	// err := pledgeActorProxy.Decline(ctx.Context())
	// if err != nil {
	// 	return input, err
	// }

	fmt.Printf("declinePledgeActivity - output: %v\n", input)
	return input, nil
}
