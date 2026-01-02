package main

import (
	"log"
	"net/http"

	"github.com/joaogeraldos/Backend-TCC/internal/database"
	"github.com/joaogeraldos/Backend-TCC/internal/player"
	ws "github.com/joaogeraldos/Backend-TCC/internal/websocket"
)

func main() {
	ws.Init()

	db := database.ConectarFirestore(ws.Ctx)
	palyerRepo := player.NewPlayerRepository(db)
	playerSvc := player.NewPlayerService(palyerRepo)
	playerHand := player.NewPlayerHandler(playerSvc)
	chatServer := ws.NewChatServer(playerSvc)

	go chatServer.RedisSubscriber(ws.Ctx)

	mux := http.NewServeMux()

	mux.HandleFunc("/chat", chatServer.ChatHandler)
	mux.HandleFunc("GET /ranking", playerHand.GetRankings)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("✅ Servidor rodando em http://localhost:8080")
	log.Println("🚀 Rota de Chat: ws://localhost:8080/chat")
	log.Println("📊 Rota de Ranking: http://localhost:8080/ranking")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("Erro ao conectar no servidor")
	}

}
