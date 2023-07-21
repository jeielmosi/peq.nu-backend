package firestore_shorten_bulk

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	config "github.com/jeielmosi/peq.nu-backend/src/config"
	"google.golang.org/api/option"
)

func getApp() (*firebase.App, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile(os.Getenv(config.FIREBASE_PATH))

	return firebase.NewApp(ctx, nil, opt)
}

func getClient() (*firestore.Client, error, context.Context) {
	ctx := context.Background()

	app, err := getApp()
	if err != nil {
		return nil, err, ctx
	}

	client, err := app.Firestore(ctx)
	return client, err, ctx
}
