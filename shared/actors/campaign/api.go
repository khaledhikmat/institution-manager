package campaignactor

import (
	"context"

	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
)

type CampaignActorStub struct {
	Id      string
	Main    func(context.Context) (campaign.Campaign, error)
	Pledges func(context.Context) ([]member.MemberPledge, error)
	Update  func(context.Context, campaign.Campaign) error
	Pledge  func(context.Context, member.MemberPledge) error
	Confirm func(context.Context, member.MemberPledge) error
	Decline func(context.Context, member.MemberPledge) error
}

func NewCampaignActor(id string) *CampaignActorStub {
	return &CampaignActorStub{Id: id}
}

func (a *CampaignActorStub) Type() string {
	return "CampaignActorType"
}

func (a *CampaignActorStub) ID() string {
	return a.Id
}
