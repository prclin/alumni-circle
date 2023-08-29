package websocket

type MessageBroker struct {
	connections []*Client
	Channel     chan Message
}

func NewMessageBroker(connections []*Client, channel chan Message) *MessageBroker {
	return &MessageBroker{connections: connections, Channel: channel}
}

func DefaultBroker() *MessageBroker {
	return &MessageBroker{connections: make([]*Client, 0), Channel: make(chan Message, 0)}
}

func (broker *MessageBroker) ProxyHandle() {
	for {
		select {
		case msg := <-broker.Channel:
			for _, connection := range broker.connections {
				connection.Conn.WriteJSON(msg)
			}
		}
	}
}

func (broker *MessageBroker) AddClient(client *Client) {
	broker.connections = append(broker.connections, client)
}
