package flow

import (
	"fmt"
	"time"

	"github.com/dapr/go-sdk/workflow"
)

func PrePledgePaymentWorkflow(ctx *workflow.WorkflowContext) (any, error) {
	var input CampaignPledge
	if err := ctx.GetInput(&input); err != nil {
		fmt.Printf("PrePledgePaymentWorkflow - error: %v\n", err)
		return nil, err
	}

	fmt.Printf("PrePledgePaymentWorkflow - input: %v\n", input)

	var output CampaignPledge

	// Wait on either timeout or paid event
	b := input.Campaign.GetBehavior()
	err := ctx.WaitForExternalEvent(PledgePaidEvent, time.Second*time.Duration(b.WaitOnPayment)).Await(nil)
	if err != nil {
		fmt.Printf("PrePledgePaymentWorkflow - error2: %v\n", err)
		// TODO: I am assuming this is a timeout condition
		// Nothing happens if we timeout....the pledge has not been submitted
		return nil, err
	}

	if err := ctx.CallActivity(submitPledgeActivity, workflow.ActivityInput(input)).Await(&input); err != nil {
		fmt.Printf("PrePledgePaymentWorkflow - error3: %v\n", err)
		return nil, err
	}

	fmt.Printf("PrePledgePaymentWorkflow - output: %v\n", input)
	return output, nil
}
