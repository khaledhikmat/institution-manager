package institutionactor

import (
	"context"

	"github.com/khaledhikmat/campaign-manager/shared/service/institution"
	"github.com/khaledhikmat/campaign-manager/shared/service/member"
)

type InstitutionActorStub struct {
	Id     string
	Main   func(context.Context) (institution.Institution, error)
	Update func(context.Context, institution.Institution) error
	Pledge func(context.Context, member.MemberPledge) error
}

func NewInstitutionActor(id string) *InstitutionActorStub {
	return &InstitutionActorStub{Id: id}
}

func (a *InstitutionActorStub) Type() string {
	return "InstitutionActorType"
}

func (a *InstitutionActorStub) ID() string {
	return a.Id
}
