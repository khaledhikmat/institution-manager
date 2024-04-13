package memberactor

import (
	"context"

	"github.com/khaledhikmat/institution-manager/shared/service/institution"
	"github.com/khaledhikmat/institution-manager/shared/service/member"
)

type MemberActorStub struct {
	Id           string
	Main         func(context.Context) (member.Member, error)
	Institutions func(context.Context) ([]institution.Institution, error)
	Memberships  func(context.Context) ([]member.MemberMembership, error)
	Transactions func(context.Context) ([]member.MemberTransaction, error)
	Pledges      func(context.Context) ([]member.MemberPledge, error)
	Update       func(context.Context, member.Member) error
	Transact     func(context.Context, member.MemberTransaction) error
	Pledge       func(context.Context, member.MemberPledge) error
}

func NewMemberActor(id string) *MemberActorStub {
	return &MemberActorStub{Id: id}
}

func (a *MemberActorStub) Type() string {
	return "MemberActorType"
}

func (a *MemberActorStub) ID() string {
	return a.Id
}
