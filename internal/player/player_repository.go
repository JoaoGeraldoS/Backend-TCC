package player

import (
	"context"
	"fmt"
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

func (r *PlayerRepository) SavePlayer(ctx context.Context, p *Player) error {
	_, err := r.client.Collection("Ordem").Doc(p.NickName).Set(ctx, p)
	if err != nil {
		return fmt.Errorf("falha ao sincronizar dados do jogador: %v", err)
	}
	return nil
}

func (r *PlayerRepository) GetRankings(ctx context.Context) ([]Player, error) {
	ranking := []Player{}

	iter := r.client.Collection("Ordem").OrderBy("pontos", firestore.Desc).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var player Player
		if err := doc.DataTo(&player); err != nil {
			log.Printf("Erro ao converter documento %s: %v", doc.Ref.ID, err)
			continue
		}
		ranking = append(ranking, player)
	}

	return ranking, nil

}
func (r *PlayerRepository) FilterName(ctx context.Context, player string) ([]Player, error) {
	var players []Player

	limiteSuperior := player + "\uf8ff"

	iter := r.client.Collection("Ordem").
		Where("usuario", ">=", player).
		Where("usuario", "<=", limiteSuperior).Documents(ctx)

	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("erro ao iterar: %v", err)
		}

		var p Player
		if err := doc.DataTo(&p); err != nil {
			return nil, fmt.Errorf("erro ao converter dados: %v", err)
		}

		players = append(players, p)
	}

	return players, nil
}

func (r *PlayerRepository) GetName(ctx context.Context, userId string) *Player {
	dsnap, err := r.client.Collection("Ordem").Doc(userId).Get(ctx)
	if err != nil {
		return nil
	}

	var p Player
	dsnap.DataTo(&p)
	if p.NickName == "" {
		return nil
	}
	return &p
}
