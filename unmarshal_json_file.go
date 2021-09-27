package golib

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

func UnmarshalJsonFile(path string, out interface{}) error {
	bin, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "Could not read file from path")
	}

	err = json.Unmarshal(bin, out)
	if err != nil {
		return errors.Wrap(err, "Could not unmarshal file")
	}

	return nil
}
