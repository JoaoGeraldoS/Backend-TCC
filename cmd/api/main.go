package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joaogeraldos/Backend-TCC/internal/database"
	"github.com/joaogeraldos/Backend-TCC/internal/middleware"
	"github.com/joaogeraldos/Backend-TCC/internal/player"
	ws "github.com/joaogeraldos/Backend-TCC/internal/websocket"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ws.Init()
	db := database.ConectarFirestore(ws.Ctx)

	mux := http.NewServeMux()
	handler := middleware.CorsMiddleware(mux)

	palyerRepo := player.NewPlayerRepository(db)
	playerSvc := player.NewPlayerService(palyerRepo)
	playerHand := player.NewPlayerHandler(playerSvc)
	chatServer := ws.NewChatServer(playerSvc)
	go chatServer.RedisSubscriber(ws.Ctx)

	mux.HandleFunc("/chat", chatServer.ChatHandler)
	mux.HandleFunc("GET /ranking", playerHand.GetRankings)
	mux.HandleFunc("GET /filter", playerHand.FilterName)
	mux.HandleFunc("POST /login", playerHand.Login)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Servidor rodando em %s", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("Erro ao conectar no servidor")
	}

}
