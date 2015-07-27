package tf2

import "github.com/james4k/rcon"

func GetStatus(host, password string) (*Status, err) {
	rc, err := rcon.Dial(host, password)
	if err != nil {
		return nil, err
	}
}
