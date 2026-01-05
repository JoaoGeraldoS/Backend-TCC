package player

import (
	"context"
)

type Player struct {
	NickName string `firestore:"usuario" json:"nick"`
	Ponts    int    `firestore:"pontos" json:"points"`
	Level    int    `firestore:"level" json:"level"`
}

type ReadPlayers interface {
	GetRankings(ctx context.Context) ([]Player, error)
	GetName(ctx context.Context, userId string) string
	FilterName(ctx context.Context, player string) ([]Player, error)
}
