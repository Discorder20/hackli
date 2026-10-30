package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

func VerifyStructure(path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	// check if it includes hackli.json
	hackliJsonPath := filepath.Join(absolutePath, "hackli.json")
	if _, err := os.Stat(hackliJsonPath); os.IsNotExist(err) {
		return fmt.Errorf("hackli.json not found in project path: %s", absolutePath)
	}
	return nil
}

func OpenProject(path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return VerifyStructure(absolutePath)
}
