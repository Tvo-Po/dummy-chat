package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

type ClientManager interface {
	Connect(client *websocket.Conn)
}

type Server struct {
	logger  *slog.Logger
	manager ClientManager
	port    string
}

func New(logger *slog.Logger, manager ClientManager, port string) *Server {
	return &Server{logger, manager, port}
}

func (s *Server) Serve() {
	http.HandleFunc("/ws", s.ws)
	http.ListenAndServe(fmt.Sprintf(":%s", s.port), nil)
	s.logger.Info("Server started...", slog.String("port", s.port))
}

func (s *Server) ws(res http.ResponseWriter, req *http.Request) {
	s.logger.Info("/ws", slog.String("method", req.Method))
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	client, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		s.logger.Error("Can't upgrade request",
			slog.String("error", err.Error()),
		)
		http.NotFound(res, req)
	}
	s.manager.Connect(client)
}
