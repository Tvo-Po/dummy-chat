package server

import (
	"context"
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
	serv    http.Server
}

func New(logger *slog.Logger, manager ClientManager, port string) *Server {
	s := &Server{logger: logger, manager: manager}
	mux := http.NewServeMux()
	mux.Handle("/ws", http.HandlerFunc(s.ws))
	s.serv = http.Server{Addr: fmt.Sprintf(":%s", port), Handler: mux}
	return s
}

func (s *Server) Serve() {
	s.logger.Info("Server started...", slog.String("port", s.serv.Addr))
	s.serv.ListenAndServe()
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

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Server releasing resources...")
	return s.serv.Shutdown(ctx)
}
