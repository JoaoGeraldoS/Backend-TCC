package player

import (
	"fmt"
	"net/http"

	"github.com/joaogeraldos/Backend-TCC/internal/middleware"
)

type PlayerHandler struct {
	svc PlayerService
}

func NewPlayerHandler(svc PlayerService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

func (h *PlayerHandler) GetRankings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	players, err := h.svc.GetRankings(ctx)
	if err != nil {
		middleware.JsonResponse(w, 500, fmt.Sprintf("Erro ao buscar playes: %v", err))
		return
	}

	middleware.JsonResponse(w, 200, players)
}
