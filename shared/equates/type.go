package equates

import (
	"github.com/khaledhikmat/institution-manager/shared/service/campaign"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
)

const (
	InstitutionManagerSecrets    = "institution-manager-secrets"
	InstitutionManagerStateStore = "institution-manager-statestore" // name must match config/redis-statestore.yaml
	InstitutionManagerPubSub     = "institution-manager-pubsub"     // name must match config/redis-pubsub.yaml
	TransactionsTopic            = "transactions-topic"
	PledgesTopic                 = "pledges-topic"
	CampaignsTopic               = "campaigns-topic"
)

type CampaignPledges struct {
	Campaign campaign.Campaign
	Pledges  []member.MemberPledge
}
