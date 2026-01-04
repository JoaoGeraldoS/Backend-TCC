package player

type PlayerRequet struct {
	Nick string `json:"nick"`
}

type PlayerResponse struct {
	TotalPoints  int      `json:"total_points"`
	TotalPlayers int      `json:"total_players"`
	OnlineNow    int      `json:"online_now"`
	Players      []Player `json:"players"`
}
