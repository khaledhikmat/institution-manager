package flow

import (
	"fmt"
	"time"

	"github.com/dapr/go-sdk/workflow"
)

func PostPledgePaymentWorkflow(ctx *workflow.WorkflowContext) (any, error) {
	var input CampaignPledge
	if err := ctx.GetInput(&input); err != nil {
		fmt.Printf("PostPledgePaymentWorkflow - error: %v\n", err)
		return nil, err
	}

	fmt.Printf("PostPledgePaymentWorkflow - input: %v\n", input)

	var output CampaignPledge

	// Wait on either timeout or paid event
	b := input.Campaign.GetBehavior()
	err := ctx.WaitForExternalEvent(PledgePaidEvent, time.Second*time.Duration(b.WaitOnPayment)).Await(nil)
	if err != nil {
		fmt.Printf("PostPledgePaymentWorkflow - error2: %v\n", err)
		// TODO: I am assuming this is a timeout condition
		// Since the pledge has already been submitted, we must decline the pledge
		if err = ctx.CallActivity(declinePledgeActivity, workflow.ActivityInput(input)).Await(&input); err != nil {
			fmt.Printf("PostPledgePaymentWorkflow - error3: %v\n", err)
			return nil, err
		}

		fmt.Printf("PostPledgePaymentWorkflow - error4: %v\n", err)
		return nil, err
	}

	if err := ctx.CallActivity(confirmPledgeActivity, workflow.ActivityInput(input)).Await(&input); err != nil {
		fmt.Printf("PostPledgePaymentWorkflow - error5: %v\n", err)
		return nil, err
	}

	fmt.Printf("PostPledgePaymentWorkflow - output: %v\n", input)
	return output, nil
}
