package player

type PlayerRequet struct {
	Nick string `json:"nick"`
}

type PlayerResponse struct {
	TotalPoints  int      `json:"total_points"`
	TotalPlayers int      `json:"total_players"`
	Players      []Player `json:"players"`
}
