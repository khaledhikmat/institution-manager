package institution

type IService interface {
	GetInstitutions(search string, opts ...InstitutionOpt) ([]Institution, error)
	GetInstitution(id string) (Institution, error)

	NewInstitution(c Institution) (Institution, error)
	UpdateInstitution(c *Institution) error

	Finalize()
}
