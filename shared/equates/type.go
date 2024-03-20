package equates

import (
	"github.com/khaledhikmat/campaign-manager/shared/service/campaign"
	"github.com/khaledhikmat/campaign-manager/shared/service/member"
)

const (
	CAMPAIGN_PUB_SUB   = "campaign-pubsub" // name must match config/redis-pubsub.yaml
	TRANSACTIONS_TOPIC = "transactions-topic"
	PLEDGES_TOPIC      = "pledges-topic"
	CAMPAIGNS_TOPIC    = "campaigns-topic"
)

type CampaignPledges struct {
	Campaign campaign.Campaign
	Pledges  []member.MemberPledge
}
