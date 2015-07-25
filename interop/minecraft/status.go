package minecraft

import (
	"strings"

	"github.com/bearbin/mcgorcon"
)

// Status is the status of the game server.
type Status struct {
	Online   bool     `json:"online"`
	Players  []string `json:"players"`
	GameName string   `json:"gamename"`
	Mode     string   `json:"mode"`
}

// Query looks up information on the remote minecraft server and returns the server Status or an error.
func Query(server string, port int, password string) (s *Status, err error) {
	s = &Status{
		GameName: "Minecraft",
		Mode:     "survival",
	}

	client, err := mcgorcon.Dial(server, port, password)
	if err != nil {
		if err.Error() == "Bad auth, could not authenticate." {
			return nil, err
		}

		s.Online = false

		return s, nil
	}

	s.Online = true

	rawPlayerList, err := client.SendCommand("list")
	if err != nil {
		return nil, err
	}

	rawPlayers := strings.Split(rawPlayerList, "online:")
	if len(rawPlayers) == 1 {
		return
	}

	players := strings.Split(rawPlayers[1], ", ")
	if players[0] != "" {
		s.Players = players
	}

	return
}
