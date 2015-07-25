package xonotic

type Status struct {
	GameName    string `json:"gamename"`
	ModName     string `json:"modname"`
	GameVersion string `json:"gameversion"`
	MaxClients  int    `json:"maxclients"`
	Clients     int    `json:"clients"`
	Bots        int    `json:"bots"`
	Mapname     string `json:"mapname"`
	Hostname    string `json:"hostname"`
	QCStatus    string `json:"qcstatus"`
}
