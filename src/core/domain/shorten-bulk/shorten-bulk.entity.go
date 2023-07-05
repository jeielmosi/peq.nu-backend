package entities

type ShortenBulkEntity struct {
	URL    string
	Clicks int64
}

func NewShortenBulkEntity(
	url string,
	clicks int64,
) *ShortenBulkEntity {
	return &ShortenBulkEntity{
		URL:    url,
		Clicks: clicks,
	}
}
