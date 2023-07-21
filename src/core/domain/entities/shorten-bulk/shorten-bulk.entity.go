package entities_shorten_bulk

type ShortenBulkEntity struct {
	URL    string `json:"url"`
	Clicks int64  `json:"clicks"`
	Custom bool   `json:"custom"`
}

func NewShortenBulkEntity(
	url string,
	clicks int64,
	custom bool,
) *ShortenBulkEntity {
	return &ShortenBulkEntity{
		URL:    url,
		Clicks: clicks,
		Custom: custom,
	}
}
