package podstalk

const Topic = "visit"

type Visit struct {
	Host       string `json:"host"`
	Method     string `json:"method"`
	Proto      string `json:"proto"`
	RemoteAddr string `json:"remote_addr"`
	Path       string `json:"path"`
	UserAgent  string `json:"user_agent"`
}
