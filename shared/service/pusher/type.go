package pusher

import (
	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
)

type IService interface {
	Externalize(campaign campaign.Campaign, pledges []member.MemberPledge) error
}
