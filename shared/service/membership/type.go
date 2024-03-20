package membership

import (
	"github.com/khaledhikmat/campaign-manager/shared/service/member"
)

type IService interface {
	GetMemberships(search string, opts ...MembershipOpt) ([]Membership, error)
	GetMembership(id string) (Membership, error)

	NewMembership(c Membership) (Membership, error)
	UpdateMembership(c *Membership) error

	GetTransactions() ([]member.MemberTransaction, error)

	Finalize()
}
