package institution

type InstitutionOpt func(i Institution)

func WithCity(t string) InstitutionOpt {
	return func(i Institution) {
		i.City = t
	}
}

func WithState(t string) InstitutionOpt {
	return func(i Institution) {
		i.State = t
	}
}

func WithCountry(t string) InstitutionOpt {
	return func(i Institution) {
		i.Country = t
	}
}

type Institution struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	Pledges     int64  `json:"pledges"`
	Donors      int64  `json:"donors"`
	Members     int64  `json:"members"`
	Memberships int64  `json:"memberships"`
}
