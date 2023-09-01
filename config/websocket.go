package config

type Websocket struct {
	Upgrader *Upgrader
}

type Upgrader struct {
	ReadBufferSize    int
	WriteBufferSize   int
	EnableCompression bool
}

var DefaultWebsocket = &Websocket{
	Upgrader: &Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: false,
	},
}
