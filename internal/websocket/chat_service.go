package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joaogeraldos/Backend-TCC/internal/player"
)

type ChatServer struct {
	Service player.PlayerService
	clients map[*websocket.Conn]bool
	mutex   sync.Mutex
}

func NewChatServer(svc player.PlayerService) *ChatServer {
	return &ChatServer{
		Service: svc,
		clients: make(map[*websocket.Conn]bool),
	}
}

func (s *ChatServer) GetConnectedCount() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.clients)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *ChatServer) ChatHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		userID = "Visitante"
	}
	nomeOficial := s.Service.GetPlayerName(r.Context(), userID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro no upgrade: ", err)
		return
	}

	log.Printf("Cliente %s conectado", nomeOficial.NickName)

	history := GetHistory()

	if len(history) > 0 {
		if err := conn.WriteJSON(history); err != nil {
			log.Println("Erro ao enviar histórico inicial:", err)
		} else {
			log.Printf("Histórico enviado para %s (%d mensagens)", nomeOficial.NickName, len(history))
		}
	}

	s.mutex.Lock()
	s.clients[conn] = true
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		delete(s.clients, conn)
		s.mutex.Unlock()
		conn.Close()
		log.Printf("Cliente %s desconectado", nomeOficial.NickName)
	}()

	for {
		var msg Message

		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Erro ao ler mensagem: ", err)
			break
		}

		loc, _ := time.LoadLocation("America/Sao_Paulo")

		msg.User = nomeOficial.NickName
		msg.Time = time.Now().In(loc).Format("15:04")

		id := SaveHistory(msg)

		payload := map[string]interface{}{
			"user":    msg.User,
			"message": msg.Message,
			"time":    msg.Time,
			"id":      id,
		}

		Publish(payload)
	}
}

func (s *ChatServer) RedisSubscriber(ctx context.Context) {
	sub := Rdb.Subscribe(ctx, Room)
	ch := sub.Channel()

	log.Println("Ouvindo mensagens do Redis no canal: ", Room)

	for pyload := range ch {
		var msg Message
		if err := json.Unmarshal([]byte(pyload.Payload), &msg); err != nil {
			continue
		}

		s.mutex.Lock()
		for client := range s.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(s.clients, client)
			}
		}
		s.mutex.Unlock()
	}
}
