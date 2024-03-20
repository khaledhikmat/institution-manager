package membership

import (
	"context"
	"fmt"

	"github.com/khaledhikmat/campaign-manager/shared/service/member"
)

// In-memory memberships database
var memberships []Membership = []Membership{
	{
		Id:            "100",
		InstitutionID: "100",
		Name:          "Unrwa",
		Description:   "Gaza",
		Benefits:      "Gaza Strip",
		Tiers: []MembershipTier{
			{
				Id:           "100",
				MembershipID: "100",
				Name:         "Platinum",
				Prefix:       "999888121",
				Benefits:     "Benefits....",
				Duration:     12,
				Currency:     "USD",
				Dependents:   3,
			},
		},
	},
}

type membershipService struct {
	CanxCtx context.Context
}

func NewService(canxCtx context.Context) IService {
	return &membershipService{
		CanxCtx: canxCtx,
	}
}

func (c *membershipService) GetMemberships(search string, opts ...MembershipOpt) ([]Membership, error) {

	filter := Membership{}
	for _, applyOpt := range opts {
		applyOpt(filter)
	}

	return memberships, nil
}

func (c *membershipService) GetMembership(id string) (Membership, error) {
	for _, mem := range memberships {
		if mem.Id == id {
			return mem, nil
		}
	}

	return Membership{}, fmt.Errorf("unable to find membership [%s]", id)
}

func (c *membershipService) NewMembership(mem Membership) (Membership, error) {
	return Membership{}, nil
}

func (c *membershipService) UpdateMembership(mem *Membership) error {
	return nil
}

func (c *membershipService) GetTransactions() ([]member.MemberTransaction, error) {
	return []member.MemberTransaction{}, nil
}

func (c *membershipService) Finalize() {
	// TODO:
}
