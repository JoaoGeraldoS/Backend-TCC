package database

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func ConectarFirestore(ctx context.Context) *firestore.Client {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	filePath := "service-account.json"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Aviso: %s não encontrado no diretório local, tentando /etc/secrets/", filePath)
		filePath = "/etc/secrets/service-account.json"
	}

	sa := option.WithCredentialsFile(filePath)

	config := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, config, sa)

	if err != nil {
		log.Fatalf("Erro ao inicializar Firebase: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Firestore: %v", err)
	}

	log.Println("✅ Firestore conectado com sucesso!")
	return client
}
