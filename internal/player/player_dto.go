package player

type PlayerRequet struct {
	Nick string `json:"nick"`
}

type PointsRequest struct {
	User   string `json:"usuario"`
	Points int    `json:"pontos"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type PlayerResponse struct {
	TotalPoints  int      `json:"total_points"`
	TotalPlayers int      `json:"total_players"`
	Players      []Player `json:"players"`
}
