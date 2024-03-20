package member

import (
	"github.com/khaledhikmat/campaign-manager/shared/service/institution"
)

type IService interface {
	GetMembers(search string, opts ...MemberOpt) ([]Member, error)
	GetMember(id string) (Member, error)

	NewMember(c Member) (Member, error)
	UpdateMember(c *Member) error

	GetInstitutions(id string) ([]institution.Institution, error)
	GetMemberships(id string) ([]MemberMembership, error)
	GetTransactions(id string) ([]MemberTransaction, error)
	GetPledges(id string) ([]MemberPledge, error)
	GetPledgesByCampaign(id string) ([]MemberPledge, error)

	// **** Institution Transactions
	Manage(memberID, institutionID string) error
	Unmanage(memberID, institutionID string) error

	// **** Membership Transactions
	NewTransaction(t MemberTransaction) error

	// **** Campaign Transactions
	NewPledge(p MemberPledge) error

	Finalize()
}
