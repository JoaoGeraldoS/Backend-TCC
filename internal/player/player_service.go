package player

import (
	"context"
	"log"
	"fmt"
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
	if p == nil {
		return fmt.Errorf("player data is nil")
	}

	getPoints := s.repo.GetName(ctx, p.NickName)
	if getPoints != nil {
		if p.Ponts < getPoints.Ponts {
			return nil
		}
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

	if s == nil || s.repo == nil {
        log.Println("Erro: Serviço de player ou repositório não inicializado")
        return &Player{NickName: "Visitante"}
    }

    if userId == "" || userId == "Visitante" {
        return &Player{NickName: "Visitante"}
    }

    player := s.repo.GetName(ctx, userId)

    if player == nil {
        return &Player{NickName: "Visitante"}
    }

    if player.NickName != userId {
        player.NickName = "Visitante"
    }

    return player
}
