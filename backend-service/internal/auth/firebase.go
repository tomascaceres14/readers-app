package auth

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

func InitFirebase() (*auth.Client, error) {
	opt := option.WithCredentialsFile("./private_keys.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	AuthClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return AuthClient, nil
}
