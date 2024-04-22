package secret

type IService interface {
	Creds(alias string) (map[string]string, error)

	Finalize()
}
