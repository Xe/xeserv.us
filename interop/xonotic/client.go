package xonotic

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Client is a rcon client for Xonotic.
type Client struct {
	host string
	port string
}

// Dial creates a client to a Xonotic server.
func Dial(host, port string) (c *Client) {
	c = &Client{
		host: host,
		port: port,
	}

	return c
}

// Status returns the server's status or an error describing the failure.
func (c *Client) Status() (stat *Status, err error) {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%s", c.host, c.port), (3 * time.Second))
	if err != nil {
		return nil, err
	}

	stat = &Status{}
	buf := make([]byte, 512)

	_, err = conn.Write([]byte("\xff\xff\xff\xffgetstatus"))
	if err != nil {
		return nil, err
	}

	_, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(buf), "\n")

	if parts[0][0] != '\xff' {
		return nil, ErrInvalidFormat
	}

	var key string

	for _, value := range strings.Split(parts[1], `\`) {
		if key == "" {
			key = value

			continue
		}

		switch key {
		case "gamename":
			stat.GameName = value
		case "modname":
			stat.ModName = value
		case "gameversion":
			stat.GameVersion = value
		case "sv_maxclients":
			stat.MaxClients, err = strconv.Atoi(value)

			if err != nil {
				return nil, err
			}
		case "clients":
			stat.Clients, err = strconv.Atoi(value)

			if err != nil {
				return nil, err
			}
		case "bots":
			stat.Bots, err = strconv.Atoi(value)

			if err != nil {
				return nil, err
			}
		case "mapname":
			stat.Mapname = value
		case "hostname":
			stat.Hostname = value
		case "qcstatus":
			stat.QCStatus = value
		}

		key = ""
	}

	return
}
