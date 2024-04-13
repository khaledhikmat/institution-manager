package flow

import (
	"fmt"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/dapr/go-sdk/workflow"
	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
)

type CampaignPledge struct {
	Campaign campaign.Campaign
	Pledge   member.MemberPledge
}

const (
	// Workflow names
	PrePledgePaymentWorkflowName  = "PrePledgePaymentWorkflow"
	PostPledgePaymentWorkflowName = "PostPledgePaymentWorkflow"

	// Activity names
	SubmitPledgeActivityName  = "SubmitPledgeActivity"
	ConfirmPledgeActivityName = "ConfirmPledgeActivity"
	DeclinePledgeActivityName = "DeclinePledgeActivity"

	// Event names
	PledgePaidEvent = "paid"
)

var workflows = map[string]workflow.Workflow{
	PrePledgePaymentWorkflowName:  PrePledgePaymentWorkflow,
	PostPledgePaymentWorkflowName: PostPledgePaymentWorkflow,
}

var activities = map[string]workflow.Activity{
	SubmitPledgeActivityName:  submitPledgeActivity,
	ConfirmPledgeActivityName: confirmPledgeActivity,
	DeclinePledgeActivityName: declinePledgeActivity,
}

// Injected DAPR client and other services
var DaprClient dapr.Client
var WorkflowClient *workflow.WorkflowWorker
var CampaignService campaign.IService
var MemberService member.IService

var WorkflowRunner *workflow.WorkflowWorker

func Start() error {
	// Create Workflow runner
	workflowRunner, err := workflow.NewWorker()
	if err != nil {
		fmt.Println("Failed to start dapr workflow engine", err)
		return err
	}

	fmt.Println("Worker initialized")

	// Register workflows
	for n, wf := range workflows {
		if err := workflowRunner.RegisterWorkflow(wf); err != nil {
			fmt.Println("Failed to register dapr workflow", err)
			return err
		}
		fmt.Printf("%s workflow registered\n", n)
	}

	// Register activities
	for n, act := range activities {
		if err := workflowRunner.RegisterActivity(act); err != nil {
			fmt.Println("Failed to register dapr workflow activity", err)
			return err
		}
		fmt.Printf("%s activity registered\n", n)
	}

	// Start workflow runner
	if err := workflowRunner.Start(); err != nil {
		fmt.Println("Failed to start dapr workflow runner", err)
		return err
	}
	fmt.Println("runner started")

	return nil
}

func Shutdown() error {
	if WorkflowRunner == nil {
		return nil
	}

	err := WorkflowRunner.Shutdown()
	if err != nil {
		return err
	}

	return nil
}
