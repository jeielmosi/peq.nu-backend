package validators

type URLBodyDto struct {
	URL string `json:"url" validate:"required,http_url"`
}
