package pusher

import (
	"github.com/khaledhikmat/campaign-manager/shared/service/campaign"
	"github.com/khaledhikmat/campaign-manager/shared/service/member"
)

type IService interface {
	Externalize(campaign campaign.Campaign, pledges []member.MemberPledge) error
}
