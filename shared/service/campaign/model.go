package campaign

import (
	"encoding/json"
)

type Opt func(c Campaign)
type Type string
type PaymentType string

const (
	// Campaign Types
	ConfirmedCampaign   Type = "Confirmed"
	UnconfirmedCampaign Type = "Unconfirmed"

	// Campaign Payment Types
	PrePledgePaymentCampaign  PaymentType = "PrePledge"
	PostPledgePaymentCampaign PaymentType = "PostPledge"
)

func WithInstitution(i string) Opt {
	return func(c Campaign) {
		c.InstitutionID = i
	}
}

func WithType(v Type) Opt {
	return func(c Campaign) {
		c.Type = v
	}
}

func WithPaymentType(v PaymentType) Opt {
	return func(c Campaign) {
		c.PaymentType = v
	}
}

func WithBehavior(v string) Opt {
	return func(c Campaign) {
		c.Behavior = v
	}
}

func WithVirtual(v bool) Opt {
	return func(c Campaign) {
		c.Virtual = v
	}
}

func WithGoal(g int64) Opt {
	return func(c Campaign) {
		c.Goal = g
	}
}

func WithCurrency(curr string) Opt {
	return func(c Campaign) {
		c.Currency = curr
	}
}

type Behavior struct {
	WaitOnPayment int64 `json:"waitOnPayment"`
}

func DefaultCampaignBehavior() Behavior {
	return Behavior{
		WaitOnPayment: 60,
	}
}

type Campaign struct {
	ID                 string      `json:"id"`
	Type               Type        `json:"type"`
	PaymentType        PaymentType `json:"paymentType"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	InstitutionID      string      `json:"institutionId"`
	Behavior           string      `json:"behavior"`
	Virtual            bool        `json:"virtual"`
	ImageURL           string      `json:"imageUrl"`
	Goal               int64       `json:"goal"`
	Currency           string      `json:"currency"`
	Duration           int64       `json:"duration"`
	GoalUSD            int64       `json:"goalUsd"`
	AwayFromGoal       int64       `json:"awayFromGoal"`
	MaxPledgeAmount    int64       `json:"maxPledgeAmount"`
	MinPledgeAmount    int64       `json:"minPledgeAmount"`
	AvgPledgeAmount    int64       `json:"avgPledgeAmount"`
	TotPledgeAmount    int64       `json:"totPledgeAmount"`
	TotPledgeAmountUSD int64       `json:"totPledgeAmountUSD"`
	Pledges            int64       `json:"pledges"`
	Donors             int64       `json:"donors"`
	// StartTime          null.Time   `json:"startTime"`
	// EndTime            null.Time   `json:"endTime"`
}

func (c Campaign) GetBehavior() Behavior {
	var b = DefaultCampaignBehavior()
	if c.Behavior != "" {
		_ = json.Unmarshal([]byte(c.Behavior), &b)
	}

	return b
}
