package membershipactor

import (
	"context"

	"github.com/khaledhikmat/institution-manager/shared/service/member"
	"github.com/khaledhikmat/institution-manager/shared/service/membership"
)

type MembershipActorStub struct {
	Id           string
	Main         func(context.Context) (member.Member, error)
	Transactions func(context.Context) ([]member.MemberTransaction, error)
	Update       func(context.Context, membership.Membership) error
	Transact     func(context.Context, member.MemberTransaction) error
}

func NewMembershipActor(id string) *MembershipActorStub {
	return &MembershipActorStub{Id: id}
}

func (a *MembershipActorStub) Type() string {
	return "MembershipActorType"
}

func (a *MembershipActorStub) ID() string {
	return a.Id
}
