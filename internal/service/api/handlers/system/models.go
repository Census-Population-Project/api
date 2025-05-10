package system

type StatusResponse struct {
	Version    string `json:"version"`
	ServerTime int64  `json:"server_time"`
	DevMode    bool   `json:"dev_mode"`
}
