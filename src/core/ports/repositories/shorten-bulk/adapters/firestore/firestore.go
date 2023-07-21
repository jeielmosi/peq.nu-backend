package firestore_shorten_bulk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	entities_shorten_bulk "github.com/jeielmosi/peq.nu-backend/src/core/domain/entities/shorten-bulk"
	repository_helpers "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/helpers"
	shorten_bulk "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	types "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/shorten-bulk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShortenBulkFirestore struct{}

func (f *ShortenBulkFirestore) Get(hash string) (
	*types.ShortenBulkRepositoryDTO,
	error,
) {
	client, err, ctx := getClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	snapshot, err := client.Collection(ShortenBulkCollection).Doc(hash).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		log.Printf("firestore.Get: %s\n", err.Error())
		return nil, err
	}

	var ans types.ShortenBulkRepositoryDTO
	err = ans.UnmarshalMap(snapshot.Data())
	if err != nil {
		return nil, err
	}

	return &ans, nil
}

func (f *ShortenBulkFirestore) GetAndIncrement(hash string, updatedAt time.Time) (
	*types.ShortenBulkRepositoryDTO,
	error,
) {
	client, err, ctx := getClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ref := client.Collection(ShortenBulkCollection).Doc(hash)
	var ans types.ShortenBulkRepositoryDTO
	getAndIncrement := func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if err != nil {
			return err
		}

		flatten := doc.Data()
		clicks := flatten[entities_shorten_bulk.ClicksField].(int64) + 1

		flatten[entities_shorten_bulk.ClicksField] = clicks
		flatten[types.UpdatedAtField] = repository_helpers.TimeToTimestamp1e8(updatedAt)
		err = ans.UnmarshalMap(flatten)
		if err != nil {
			log.Println("GetAndIncrement Error:", err.Error())
			return err
		}
		return tx.Set(ref, flatten, firestore.MergeAll)
	}

	err = client.RunTransaction(ctx, getAndIncrement)
	if err != nil {
		return nil, err
	}

	return &ans, nil
}

func (f *ShortenBulkFirestore) GetOldest(size int) (
	map[string]*types.ShortenBulkRepositoryDTO,
	error,
) {
	client, err, ctx := getClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	iter := client.
		Collection(ShortenBulkCollection).
		OrderBy(types.CreatedAtField, firestore.Asc).
		Where(entities_shorten_bulk.CustomField, "!=", true).
		Limit(size).
		Documents(ctx)

	mp := map[string]*types.ShortenBulkRepositoryDTO{}
	for true {
		snapshot, err := iter.Next()
		if err != nil {
			break
		}

		var dto types.ShortenBulkRepositoryDTO
		err = dto.UnmarshalMap(snapshot.Data())
		if err != nil {
			log.Println(err.Error())
			continue
		}

		mp[snapshot.Ref.ID] = &dto
	}

	return mp, nil
}

func (f *ShortenBulkFirestore) PostSafe(
	hash string,
	dto types.ShortenBulkRepositoryDTO,
) error {
	client, err, ctx := getClient()
	if err != nil {
		return err
	}
	defer client.Close()

	flatten, err := dto.MarshalMap()
	ref := client.Collection(ShortenBulkCollection).Doc(hash)
	postSafe := func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := tx.Get(ref)
		if err == nil {
			return errors.New(fmt.Sprintf("Hash '%s' is already used", hash))
		}

		if status.Code(err) != codes.NotFound {
			return err
		}
		return tx.Set(ref, flatten, firestore.MergeAll)
	}

	return client.RunTransaction(ctx, postSafe)
}

func (f *ShortenBulkFirestore) PostUnsafe(
	hash string,
	dto types.ShortenBulkRepositoryDTO,
) error {
	client, err, ctx := getClient()
	if err != nil {
		return err
	}
	defer client.Close()

	flatten, err := dto.MarshalMap()
	ref := client.Collection(ShortenBulkCollection).Doc(hash)
	postUnsafe := func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}

		custom_interface, err := doc.DataAt(entities_shorten_bulk.CustomField)
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}

		custom, ok := custom_interface.(bool)
		if ok && custom {
			return errors.New("Not allowed to overwrite a custom hash")
		}

		return tx.Set(ref, flatten, firestore.MergeAll)
	}

	return client.RunTransaction(ctx, postUnsafe)
}

func NewShortenBulkFirestore() shorten_bulk.ShortenBulkRepository {
	return &ShortenBulkFirestore{}
}
