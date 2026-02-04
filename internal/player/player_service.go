package player

import (
	"context"
	"log"
)

type PlayerService interface {
	GetRankings(ctx context.Context) ([]Player, error)
	GetPlayerName(ctx context.Context, id string) *Player
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
	getPoints := s.repo.GetName(ctx, p.NickName)

	if p.Ponts < getPoints.Ponts {
		return nil
	}

	return s.repo.SavePlayer(ctx, p)
}

func (s *servicePlayer) GetRankings(ctx context.Context) ([]Player, error) {
	return s.repo.GetRankings(ctx)
}

func (s *servicePlayer) FilterName(ctx context.Context, player string) ([]Player, error) {

	return s.repo.FilterName(ctx, player)
}

func (s *servicePlayer) GetPlayerName(ctx context.Context, userId string) *Player {

	player := s.repo.GetName(ctx, userId)

	if userId == "" || userId == "Visitante" {
		log.Println("Busca abortada: ID inválido ou reservado")
		player.NickName = "Visitante"
	}

	if player.NickName != userId {
		player.NickName = "Visitante"
	}

	return player
}
