package manager

import (
	"dummy-chat/internal/domain"
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
)

type ClientManager struct {
	logger       *slog.Logger
	clients      map[*websocket.Conn]string
	broadcast    chan domain.Message
	disconnected chan *websocket.Conn
}

func New(logger *slog.Logger) *ClientManager {
	return &ClientManager{
		logger:       logger,
		clients:      map[*websocket.Conn]string{},
		broadcast:    make(chan domain.Message),
		disconnected: make(chan *websocket.Conn),
	}
}

func getClientLogInfo(client *websocket.Conn, name string) slog.Attr {
	return slog.Group("client",
		slog.String("addr", fmt.Sprintf("%p", client)),
		slog.String("name", name),
	)
}

func (m *ClientManager) Run() {
	m.logger.Info("Manager accepting connections...")
	for {
		select {
		// TODO: ctx done
		case client := <-m.disconnected:
			if name, ok := m.clients[client]; ok {
				m.logger.Info("Disconnecting client",
					getClientLogInfo(client, name),
				)
				client.WriteMessage(websocket.CloseMessage, []byte{})
				delete(m.clients, client)
				go func() {
					m.broadcast <- domain.Message{name, "disconnected"}
				}()
			}
		case msg := <-m.broadcast:
			m.logger.Debug(fmt.Sprintf("Broadcasting message: '%s'", msg.Content))
			for client, name := range m.clients {
				if err := client.WriteJSON(msg); err != nil {
					m.logger.Error("Can't send message to client",
						getClientLogInfo(client, name),
						slog.String("error", err.Error()),
					)
					go func() {
						m.disconnected <- client
					}()
				}
			}
		}
	}
}

func (m *ClientManager) Connect(client *websocket.Conn) {
	m.logger.Info("New websocket connection", getClientLogInfo(client, ""))
	m.clients[client] = ""
	go func() {
		_, encodedName, err := client.ReadMessage()
		if err != nil {
			m.logger.Error("Can't read client",
				getClientLogInfo(client, ""),
				slog.String("error", err.Error()),
			)
			select {
			case m.disconnected <- client:
			default:
			}
			return
		}
		name := string(encodedName)
		m.clients[client] = name
		m.logger.Info("Set client name", getClientLogInfo(client, name))
		m.broadcast <- domain.Message{name, "connected"}
		for {
			msg := domain.Message{}
			if err := client.ReadJSON(&msg); err != nil {
				m.logger.Error("Can't read client message",
					getClientLogInfo(client, name),
					slog.String("error", err.Error()),
				)
				select {
				case m.disconnected <- client:
				default:
				}
				return
			}
			m.logger.Debug(fmt.Sprintf("New message from client: '%s'", msg.Content),
				getClientLogInfo(client, name),
			)
			m.broadcast <- msg
		}
	}()
}
