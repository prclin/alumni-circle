package messaging

type Context struct {
	broker *StompBroker
	Frame  *Frame
	Conn   *Conn
	Params map[string]string
}
