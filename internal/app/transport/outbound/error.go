package outbound

type Error struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Params  interface{} `json:"params"`
}
