package secretkey

import (
	"net/http"
)

func secretKeyPageHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, "secretkey.page.tmpl", nil)
}

func secretKeyHandler(w http.ResponseWriter, r *http.Request) {
	sk := r.FormValue("secretkey")
	if sk == "" {
		data := templateData{
			Error:"secretkey cannot be empty",
		}
		render(w, r, "secretkey.page.tmpl", &data)
		return
	}
	
	err := DBModel.Update(sk)
	if err != nil {
		serverError(w, err)
		return
	}

	data := templateData{
		Success: "the secretkey has been updated\n sign it into all of your devices",
	}
	render(w, r, "secretkey.page.tmpl", &data)
}