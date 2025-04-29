package secretkey

import (
	"html/template"
	"log"
	"maps"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mohamidsaiid/uniclipboard/internal/models"
)

var templateCache = map[string]*template.Template{}
var DBModel *models.UsersModel
func StartSecertKeyWebServer(secretKeyPort string, dbmodel *models.UsersModel) {
	tc, err := newTemplateCache("./ui/html")
	maps.Copy(templateCache, tc)
	if err != nil {
		log.Println(err)
		return 
	}
	DBModel = dbmodel

	router := httprouter.New()	

	router.Handler("GET", "/static/*filepath", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.HandlerFunc(http.MethodGet, "/secretkey", secretKeyPageHandler)
	router.HandlerFunc(http.MethodPost, "/secretkey", secretKeyHandler)
	srv := &http.Server{
		Addr: secretKeyPort,
		Handler : router,
	}
	log.Println(srv.ListenAndServe())
} 

