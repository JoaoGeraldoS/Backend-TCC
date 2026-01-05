package player

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/joaogeraldos/Backend-TCC/internal/middleware"
)

type PlayerHandler struct {
	svc PlayerService
}

func NewPlayerHandler(svc PlayerService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

func (h *PlayerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var dto PlayerRequet

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		middleware.JsonResponse(w, 400, err)
		return
	}

	player := h.svc.GetPlayerName(r.Context(), dto.Nick)

	if player == "Visitante" {
		middleware.JsonResponse(w, 404, "Usuario nao existe")
		return
	}

	middleware.JsonResponse(w, 200, player)
}

func (h *PlayerHandler) FilterName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	player := r.URL.Query().Get("player")

	players, err := h.svc.FilterName(ctx, player)
	if err != nil {
		middleware.JsonResponse(w, 500, fmt.Sprintf("Erro ao buscar playes: %v", err))
		return
	}

	middleware.JsonResponse(w, 200, players)
}

func (h *PlayerHandler) GetRankings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	players, err := h.svc.GetRankings(ctx)
	if err != nil {
		middleware.JsonResponse(w, 500, fmt.Sprintf("Erro ao buscar playes: %v", err))
		return
	}

	total := 0
	for _, p := range players {
		total += p.Ponts
	}

	res := PlayerResponse{
		TotalPoints:  total,
		TotalPlayers: len(players),
		Players:      players,
	}

	middleware.JsonResponse(w, 200, res)
}
