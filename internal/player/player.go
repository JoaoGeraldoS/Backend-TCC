package player

import (
	"context"
)

type Player struct {
	NickName string `firestore:"usuario" json:"nick"`
	Ponts    int    `firestore:"pontos" json:"points"`
}

type ReadPlayers interface {
	GetRankings(ctx context.Context) ([]Player, error)
	GetName(ctx context.Context, userId string) *Player
	FilterName(ctx context.Context, player string) ([]Player, error)
}

type CreatorPlayers interface {
	SavePlayer(ctx context.Context, p *Player) error
}

type PlayerInterface interface {
	ReadPlayers
	CreatorPlayers
}
