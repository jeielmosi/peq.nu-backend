package shorten_bulk_gateway

import (
	"log"
	"reflect"
	"testing"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
)

func TestGet(
	hash string,
	gateway ShortenBulkGateway,
	t *testing.T,
) {
	before, err := gateway.Get(hash)
	if err != nil {
		t.Errorf("Get element at hash '%s': %s", hash, err.Error())
	}

	if before.Clicks%2 != 1 {
		t.Errorf("Get element at hash '%s': Count clicks is %d", hash, before.Clicks)
	}

	after, err := gateway.Get(hash)
	if err != nil {
		t.Errorf("Get element at hash '%s': %s", hash, err.Error())
	}

	if after.Clicks%2 != 0 {
		t.Errorf("Get element at hash '%s': Count clicks was %d and is %d", hash, before.Clicks, after.Clicks)
	}

	before.Clicks += 1

	if !reflect.DeepEqual(before, after) {
		t.Errorf("Result and expected are different")
	}
}

func TestPost(
	gateway ShortenBulkGateway,
	exp entities.ShortenBulkEntity,
	t *testing.T,
) {
	hash, err := gateway.Post(exp)
	if err != nil {
		t.Errorf("Test error at post entity: %s", err.Error())
		return
	}

	res, err := gateway.Get(hash)
	if err != nil {
		t.Errorf("Test error at get entity at hash '%s': %s", hash, err.Error())
		return
	}

	exp.Clicks += 1

	if !reflect.DeepEqual(res, &exp) {
		log.Println("res:", res)
		log.Println("exp:", exp)

		t.Errorf("Result and expected are different")
	}
}
