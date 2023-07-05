package firestore_shorten_bulk

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func getApp(envName string) (*firebase.App, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile(os.Getenv(envName + "_FIREBASE_PATH"))

	return firebase.NewApp(ctx, nil, opt)
}

func getClient(envName string) (*firestore.Client, error, context.Context) {
	ctx := context.Background()

	app, err := getApp(envName)
	if err != nil {
		return nil, err, ctx
	}

	client, err := app.Firestore(ctx)
	return client, err, ctx
}
