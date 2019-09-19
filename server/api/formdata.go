package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// ConvertFormData converts form data in the request to the json.
func convertFormData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			next.ServeHTTP(w, r)
			return
		}

		if err := r.ParseMultipartForm(r.ContentLength); err != nil {
			logrus.Errorf("failed to parse multipart form: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		toMarshal := make(map[string]interface{}, len(r.PostForm))
		for k, vv := range r.PostForm {
			if len(vv) > 0 {
				toMarshal[k] = vv[0]
			}
		}

		encodedJSON, err := json.Marshal(toMarshal)
		if err != nil {
			logrus.Errorf("failed to encode multipart to json: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(encodedJSON))
		r.Header.Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
