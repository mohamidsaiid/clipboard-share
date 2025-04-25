package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (srvr *Server) routes() *httprouter.Router{
	router := httprouter.New()
	
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", srvr.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/clipboard", srvr.clipboardHandler)	
	router.HandlerFunc(http.MethodGet, "/api/v1/clipboarddata", srvr.lastCopiedData)
	go srvr.handleMessages()
	return router
}