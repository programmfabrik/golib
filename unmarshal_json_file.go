package golib

import (
	"encoding/json"
	"fmt"
	"os"
)

func UnmarshalJsonFile(path string, out interface{}) error {
	bin, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Could not read file from path: %w", err)
	}

	err = json.Unmarshal(bin, out)
	if err != nil {
		return fmt.Errorf("Could not unmarshal file: %w", err)
	}

	return nil
}
