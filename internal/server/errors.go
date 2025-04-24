package server

import (
	"net/http"

	"github.com/mohamidsaiid/uniclipboard/internal/jsonParser"
)

func (srvr *Server) logError(r *http.Request, err error) {
	srvr.Logger.Println(err)
}

func (srvr *Server) errorResponse (w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	err := jsonParser.WriteJSON(w, status, message, nil)
	if err != nil {
		srvr.logError(r, err)
		w.WriteHeader(500)
	}
}

func (srvr *Server) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	srvr.logError(r, err)
	message := "The server encountered a problem and couldn't process your request"
	
	srvr.errorResponse(w, r, http.StatusInternalServerError, message)
}