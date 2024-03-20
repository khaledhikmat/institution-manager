package membership

type MembershipOpt func(m Membership)

func WithInstitution(i string) MembershipOpt {
	return func(m Membership) {
		m.InstitutionID = i
	}
}

type Membership struct {
	Id            string           `json:"id"`
	InstitutionID string           `json:"institutionId"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Benefits      string           `json:"benefits"`
	Transactions  int              `json:"transactions"`
	Purchases     int              `json:"purchases"`
	Renewals      int              `json:"renewals"`
	Cancellations int              `json:"cancellations"`
	Tiers         []MembershipTier `json:"tiers"`
}

type MembershipTier struct {
	Id           string `json:"id"`
	MembershipID string `json:"membershipId"`
	Name         string `json:"name"`
	Prefix       string `json:"prefix"`
	Benefits     string `json:"benefits"`
	Duration     int    `json:"duration"`
	Cost         int    `json:"cost"`
	Currency     string `json:"currency"`
	Dependents   int    `json:"dependents"`
}
