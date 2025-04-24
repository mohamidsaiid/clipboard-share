package jsonParser

import (
	"net/http"
	"encoding/json"
	"maps"
)

func WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// here to parse the data from its own structure to a JSON output the write
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	// adding any needed http header to the response
	maps.Copy(w.Header(), headers)

	// setting the content-type header to JSON insted of application/text
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

