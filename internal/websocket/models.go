package websocket

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
	Time    string `json:"time"`
}
