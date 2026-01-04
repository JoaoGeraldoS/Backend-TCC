package player

import (
	"context"
	"fmt"
	"log"
)

type PlayerService interface {
	GetRankings(ctx context.Context) ([]Player, error)
	GetPlayerName(ctx context.Context, id string) string
}

type servicePlayer struct {
	repo ReadPlayers
}

func NewPlayerService(repo ReadPlayers) *servicePlayer {
	return &servicePlayer{repo: repo}
}

func (s *servicePlayer) GetRankings(ctx context.Context) ([]Player, error) {
	return s.repo.GetRankings(ctx)
}

func (s *servicePlayer) GetPlayerName(ctx context.Context, id string) string {

	player := s.repo.GetName(ctx, id)

	if player != id {
		log.Println("user nao exits")
		return "Visitante"
	}

	fmt.Println("Esta vindo", player)

	return player
}
