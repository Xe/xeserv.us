package tf2

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/anmitsu/go-shlex"
	"github.com/james4k/rcon"
)

func GetStatus(host, password string) (*Status, error) {
	rc, err := rcon.Dial(host, password)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	id, err := rc.Write("status")
	if err != nil {
		return nil, err
	}

	reply, theirID, err := rc.Read()
	if err != nil {
		return nil, err
	}

	if id != theirID {
		return nil, errors.New("tf2: things got out of order somehow, I'm gonna stop trying.")
	}

	s := &Status{
		Online: true,
	}

	lines := strings.Split(reply, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		log.Printf("%q", fields)

		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "hostname:":
			s.Hostname = strings.Join(fields[1:], " ")

		case "version":
			s.Version = strings.Join(fields[2:], " ")

		case "map":
			s.MapName = fields[2]

		case "tags":
			s.Tags = fields[2]

		case "players":
			var err error

			s.HumanPlayers, err = strconv.Atoi(fields[2])
			if err != nil {
				return nil, err
			}

			s.Bots, err = strconv.Atoi(fields[4])
			if err != nil {
				return nil, err
			}

			s.MaxPlayers, err = strconv.Atoi(fields[6][1:])
			if err != nil {
				return nil, err
			}

		case "#":
			if fields[1] == "userid" {
				continue
			}

			info, err := shlex.Split(line, true)
			if err != nil {
				return nil, err
			}

			p := Player{
				UserID:    info[1],
				Name:      info[2],
				UniqueID:  info[3],
				Connected: info[4],
				State:     info[7],
				Address:   info[8],
				IsBot:     info[3] == "BOT",
			}

			s.Players = append(s.Players, p)
		}
	}

	return s, nil
}
