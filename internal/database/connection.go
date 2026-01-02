package database

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func ConectarFirestore(ctx context.Context) *firestore.Client {

	sa := option.WithCredentialsFile("service-account.json")

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("Erro ao inicializar Firebase: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Firestore: %v", err)
	}

	return client
}
