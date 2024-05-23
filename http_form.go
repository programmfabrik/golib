package golib

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
)

// ParseForm works like http.Request.ParseForm but additionally
// adds values received in json format to request.Form
func ParseForm(req *http.Request) error {

	err := req.ParseForm()
	if err != nil {
		return err
	}

	ct := req.Header.Get("content-type")

	if ct == "" { // Assume "GET"
		return nil
	}

	mediatype, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return err
	}

	switch mediatype {

	case "application/x-www-form-urlencoded":
		return nil
	case "multipart/form-data":
		err = req.ParseMultipartForm(4096)
		if err != nil {
			return fmt.Errorf("Unable to parse form: %w", err)
		}
	case "application/json":
		jsonForm := map[string]interface{}{}
		err = json.NewDecoder(req.Body).Decode(&jsonForm)
		if err != nil {
			return fmt.Errorf("ParseForm failed: %w", err)
		}
		for k, v := range jsonForm {
			switch v1 := v.(type) {
			case string:
				req.Form.Set(k, v1)
			case *string:
				req.Form.Set(k, *v1)
			default:
				bytes, _ := json.Marshal(v)
				req.Form.Set(k, string(bytes))
			}
		}
	}

	return nil
}
