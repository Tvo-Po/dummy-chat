package manager

import (
  "dummy-chat/internal/domain"
  "github.com/gorilla/websocket"
)

type ClientManager struct {
  clients map[*websocket.Conn]string
  broadcast chan domain.Message
  disconnected chan *websocket.Conn
}

func New() *ClientManager {
  return &ClientManager{
    clients: map[*websocket.Conn]string{},
    broadcast: make(chan domain.Message),
    disconnected: make(chan *websocket.Conn),
  }
}

func (m *ClientManager) Run() {
  for {
    select {
    // TODO: ctx done
    case client := <- m.disconnected:
      if name, ok := m.clients[client]; ok {
        client.WriteMessage(websocket.CloseMessage, []byte{})
        delete(m.clients, client)
        go func() {
          m.broadcast <- domain.Message{name, "disconnected"}
        }()
      }
    case msg := <- m.broadcast:
      for client, _ := range m.clients {
        if err := client.WriteJSON(msg); err != nil {
          go func() {
            m.disconnected <- client
          }()
        }
      }
    }
  }
}

func (m *ClientManager) Connect(client *websocket.Conn) {
  m.clients[client] = ""
  go func() {
    _, encodedName, err := client.ReadMessage()
    if err != nil {
      select {
      case m.disconnected <- client:
      default:
      }
      return
    }
    name := string(encodedName)
    m.clients[client] = name
    m.broadcast <- domain.Message{name, "connected"}
    for {
      msg := domain.Message{}
      if err := client.ReadJSON(msg); err != nil {
        select {
        case m.disconnected <- client:
        default:
        }
        return
      }
      m.broadcast <- msg
    }
  }()
}
