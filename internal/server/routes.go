package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (srvr *Server) routes() *httprouter.Router{
	router := httprouter.New()
	
	router.HandlerFunc(http.MethodGet, "/healthcheck", srvr.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/clipboard", srvr.clipboardHandler)	
	router.HandlerFunc(http.MethodGet, "/clipboarddata", srvr.lastCopiedData)
	go srvr.handleMessages()
	return router
}