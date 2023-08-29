package websocket

type Message struct {
	Type        string `json:"type"`
	Destination uint64 `json:"destination"`
	Body        string `json:"body"`
}

const (
	TypeUnicast   = "unicast"
	TypeBroadcast = "broadcast"
)
