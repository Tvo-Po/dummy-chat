package server

import (
  "net/http"
  "github.com/gorilla/websocket"
)

type ClientManager interface {
  Connect(client *websocket.Conn)
}

type Server struct {
  manager ClientManager
  port string
}

func New(manager ClientManager, port string) *Server {
  return &Server{manager, port}
}

func (s *Server) Serve() {
  http.HandleFunc("/ws", s.ws)
  http.ListenAndServe(s.port, nil)
}

func (s *Server) ws(res http.ResponseWriter, req *http.Request) {
  upgrader := &websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
      return true
    },
  }
  client, err := upgrader.Upgrade(res, req, nil)
  if err != nil {
    http.NotFound(res, req)
  }
  s.manager.Connect(client)
}
