package secretkey

import (
	"log"
	"net/http"
)
func serverError(w http.ResponseWriter, err error) {
	log.Println("the secretkey server has an error ",  err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
