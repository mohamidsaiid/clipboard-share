package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	uniclipboard "github.com/mohamidsaiid/uniclipboard/internal/clipboard"
	"github.com/mohamidsaiid/uniclipboard/internal/models"
)

type config struct {
	Port string
	Env  string
}

type Server struct {
	Clients   map[*websocket.Conn]bool
	Cfg       config
	Logger    *log.Logger
	clipboard *uniclipboard.UniClipboard
	user      *models.UsersModel
}

func NewServer(port string, clipboard *uniclipboard.UniClipboard, u *models.UsersModel) *Server {
	return &Server{
		Clients: make(map[*websocket.Conn]bool),
		Cfg: config{
			Port: port,
			Env:  "development",
		},
		Logger:    log.New(os.Stdout, "", log.Ldate|log.Ltime),
		clipboard: clipboard,
		user : u,
	}
}

func (s *Server) Start() {
	srv := &http.Server{
		Addr:    s.Cfg.Port,
		Handler: s.routes(),
	}

	s.Logger.Printf("starting server %s on port %s", s.Cfg.Env, srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
