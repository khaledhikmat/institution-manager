package member

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"

	"github.com/khaledhikmat/institution-manager/shared/service/institution"
)

type MemberOpt func(m Member)

func WithType(t string) MemberOpt {
	return func(m Member) {
		m.Type = t
	}
}

func WithRole(r string) MemberOpt {
	return func(m Member) {
		m.Role = r
	}
}

type MemberRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MemberType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Member struct {
	ID            string                    `json:"id"`
	Role          string                    `json:"role"`
	Type          string                    `json:"type"`
	Parent        string                    `json:"parent"`
	ExternalID    string                    `json:"externalId"`
	Name          string                    `json:"name"`
	Email         string                    `json:"email"`
	Phone         string                    `json:"phone"`
	Memberships   int                       `json:"memberships"`
	Dependents    int                       `json:"dependents"`
	Pledges       int                       `json:"pledges"`
	Transactions  int                       `json:"transactions"`
	Purchases     int                       `json:"purchases"`
	Renewals      int                       `json:"renewals"`
	Cancellations int                       `json:"cancellations"`
	Institutions  []institution.Institution `json:"institutions"`
}

type MemberMembership struct {
	ID           string    `json:"id"`
	MemberID     string    `json:"memberId"`
	MembershipID string    `json:"membershipId"`
	Tier         string    `json:"tier"`
	Number       string    `json:"number"`
	SinceDate    time.Time `json:"sinceDate"`
	ExpDate      time.Time `json:"expDate"`
	Renewals     int64     `json:"renewals"`
	Dependents   []Member  `json:"dependents"`
}

type MemberTransaction struct {
	ID                  string    `json:"id"`
	MemberID            string    `json:"memberId"`
	MembershipID        string    `json:"membershipId"`
	Tier                string    `json:"tier"`
	Type                string    `json:"type"`
	StartDate           time.Time `json:"startDate"`
	ExpDate             time.Time `json:"expDate"`
	Amount              int64     `json:"amount"`
	ExgRate             float64   `json:"exgRate"`
	PaymentMethod       string    `json:"paymentMethod"`
	PaymentConfirmation string    `json:"paymentConfirmation"`
}

type MemberPledge struct {
	ID            string    `json:"id"`
	MemberID      string    `json:"memberId"`
	CampaignID    string    `json:"campaignId"`
	Time          time.Time `json:"time"`
	Amount        int64     `json:"amount"`
	ExgRate       float64   `json:"exgRate"`
	ConfirmedTime null.Time `json:"confirmedTime"`
	PersistedTime null.Time `json:"persistedTime"`
}

func NewMemberPledge() MemberPledge {
	p := MemberPledge{}
	p.ID = uuid.New().String() // Generate a unique ID

	return p
}
