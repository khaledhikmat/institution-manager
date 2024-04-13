package campaign

type IService interface {
	GetCampaigns(search string, opts ...Opt) ([]Campaign, error)
	GetCampaign(id string) (Campaign, error)

	NewCampaign(c Campaign) (Campaign, error)
	UpdateCampaign(c *Campaign) error

	Finalize()
}
