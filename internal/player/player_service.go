package player

import (
	"context"
	"log"
)

type PlayerService interface {
	GetRankings(ctx context.Context) ([]Player, error)
	GetPlayerName(ctx context.Context, id string) string
	FilterName(ctx context.Context, player string) ([]Player, error)
	SavePlayer(ctx context.Context, p *Player) error
}

type servicePlayer struct {
	repo PlayerInterface
}

func NewPlayerService(repo PlayerInterface) *servicePlayer {
	return &servicePlayer{repo: repo}
}

func (s *servicePlayer) SavePlayer(ctx context.Context, p *Player) error {
	return s.repo.SavePlayer(ctx, p)
}

func (s *servicePlayer) GetRankings(ctx context.Context) ([]Player, error) {
	return s.repo.GetRankings(ctx)
}

func (s *servicePlayer) FilterName(ctx context.Context, player string) ([]Player, error) {

	return s.repo.FilterName(ctx, player)
}

func (s *servicePlayer) GetPlayerName(ctx context.Context, userId string) string {

	player := s.repo.GetName(ctx, userId)

	if userId == "" || userId == "Visitante" {
		log.Println("Busca abortada: ID inválido ou reservado")
		return "Visitante"
	}

	if player != userId {
		return "Visitante"
	}

	return player
}
