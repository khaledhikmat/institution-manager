package pusher

import (
	"context"

	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
)

type pusherService struct {
	CanxCtx context.Context
}

func NewService(canxCtx context.Context) IService {
	return &pusherService{
		CanxCtx: canxCtx,
	}
}

func (c *pusherService) Externalize(campaign campaign.Campaign, pledges []member.MemberPledge) error {
	return nil
}
