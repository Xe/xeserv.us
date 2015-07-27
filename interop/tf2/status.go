package tf2

// Status is the status of the game server.
type Status struct {
	Online       bool     `json:"online"`
	Hostname     string   `json:"hostname"`
	Version      string   `json:"version"`
	Tags         string   `json:"tags"`
	Players      []Player `json:"players"`
	HumanPlayers int      `json:"human_players"`
	Bots         int      `json:"bots"`
	MaxPlayers   int      `json:"max_players"`
	MapName      string   `json:"map_name"`
}

// Player represents an individual player connected to the server.
type Player struct {
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	UniqueID  string `json:"unique_id"`
	Connected string `json:"connected"`
	State     string `json:"state"`
	Address   string `json:"addr"`
	IsBot     bool   `json:"bot"`
}
