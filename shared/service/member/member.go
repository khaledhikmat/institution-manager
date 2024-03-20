package member

import (
	"context"
	"fmt"

	"github.com/khaledhikmat/campaign-manager/shared/service/institution"
)

// In-memory members database
var members []Member = []Member{
	{
		Id:         "100",
		Role:       "Admin",
		Type:       "Donor",
		ExternalID: "1000-9993-6767-2387-9000",
		Name:       "Omar Fateh",
		Email:      "omar.fateh@gmail.com",
		Phone:      "210-555-1212",
	},
	{
		Id:         "200",
		Role:       "Manager",
		Type:       "Donor",
		ExternalID: "1000-9993-6767-2387-8000",
		Name:       "Shk Yousuf",
		Email:      "shk.yousuf@gmail.com",
		Phone:      "210-555-1212",
		Institutions: []institution.Institution{
			{
				Id:      "100",
				Name:    "Unrwa",
				City:    "Gaza",
				State:   "Gaza Strip",
				Country: "Palestine",
			},
		},
	},
}

type memberService struct {
	CanxCtx context.Context
}

func NewService(canxCtx context.Context) IService {
	return &memberService{
		CanxCtx: canxCtx,
	}
}

func (c *memberService) GetMembers(search string, opts ...MemberOpt) ([]Member, error) {

	filter := Member{}
	for _, applyOpt := range opts {
		applyOpt(filter)
	}

	return members, nil
}

func (c *memberService) GetMember(id string) (Member, error) {
	for _, mem := range members {
		if mem.Id == id {
			return mem, nil
		}
	}

	return Member{}, fmt.Errorf("unable to find member [%s]", id)
}

func (c *memberService) NewMember(mem Member) (Member, error) {
	return Member{}, nil
}

func (c *memberService) UpdateMember(mem *Member) error {
	return nil
}

func (c *memberService) GetInstitutions(id string) ([]institution.Institution, error) {
	return []institution.Institution{}, nil
}

func (c *memberService) GetMemberships(id string) ([]MemberMembership, error) {
	return []MemberMembership{}, nil
}

func (c *memberService) GetTransactions(id string) ([]MemberTransaction, error) {
	return []MemberTransaction{}, nil
}

func (c *memberService) GetPledges(id string) ([]MemberPledge, error) {
	return []MemberPledge{}, nil
}

func (c *memberService) GetPledgesByCampaign(id string) ([]MemberPledge, error) {
	return []MemberPledge{}, nil
}

// **** Institution Transactions
func (c *memberService) Manage(memberID, institutionID string) error {
	return nil
}

func (c *memberService) Unmanage(memberID, institutionID string) error {
	return nil
}

// **** Membership Transactions
func (c *memberService) NewTransaction(t MemberTransaction) error {
	return nil
}

// **** Campaign Transactions
func (c *memberService) NewPledge(p MemberPledge) error {
	return nil
}

func (c *memberService) Finalize() {
	// TODO:
}
