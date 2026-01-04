package player

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type PlayerRepository struct {
	client *firestore.Client
}

func NewPlayerRepository(cli *firestore.Client) *PlayerRepository {
	return &PlayerRepository{client: cli}
}

func (r *PlayerRepository) GetRankings(ctx context.Context) ([]Player, error) {
	client := r.client

	var ranking []Player

	iter := client.Collection("Ordem").OrderBy("pontos", firestore.Desc).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var player Player
		doc.DataTo(&player)
		ranking = append(ranking, player)
	}

	return ranking, nil

}

func (r *PlayerRepository) GetName(ctx context.Context, userId string) string {

	if userId == "" || userId == "Visitante" {
		log.Println("Busca abortada: ID inválido ou reservado")
		return "Visitante"
	}
	dsnap, err := r.client.Collection("Ordem").Doc(userId).Get(ctx)
	if err != nil {
		log.Printf("Erro no Firestore para o ID %s: %v", userId, err)
		return "Visitante"
	}

	nome := dsnap.Data()["usuario"].(string)
	log.Printf("Sucesso! Nome recuperado: %s", nome)
	return nome
}
