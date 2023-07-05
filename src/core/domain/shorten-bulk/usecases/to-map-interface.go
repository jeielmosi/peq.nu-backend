package usecases

import entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"

func ToMapInterface(shortenBulk *entities.ShortenBulkEntity) map[string]interface{} {
	ans := map[string]interface{}{}

	if shortenBulk == nil {
		return ans
	}

	ans["url"] = shortenBulk.URL
	ans["clicks"] = shortenBulk.Clicks

	return ans
}
