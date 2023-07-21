package entities_shorten_bulk

import (
	"errors"
	"fmt"
)

const (
	URLField    = "url"
	ClicksField = "clicks"
	CustomField = "custom"
)

func (s *ShortenBulkEntity) MarshalMap() (map[string]interface{}, error) {
	if s == nil {
		return nil, errors.New("ShortenBulkEntity is nil")
	}

	ans := make(map[string]interface{})
	ans[URLField] = s.URL
	ans[ClicksField] = s.Clicks
	ans[CustomField] = s.Custom

	return ans, nil
}

func (s *ShortenBulkEntity) UnmarshalMap(mp map[string]interface{}) error {
	template := "The key '%s' not found on map[string]interface{}"

	clicks, ok := mp[ClicksField].(int64)
	if !ok {
		return errors.New(fmt.Sprintf(template, ClicksField))
	}

	url, ok := mp[URLField].(string)
	if !ok {
		return errors.New(fmt.Sprintf(template, URLField))
	}

	custom, ok := mp[CustomField].(bool)
	if !ok {
		custom = false
	}

	*s = *NewShortenBulkEntity(url, clicks, custom)

	return nil
}
