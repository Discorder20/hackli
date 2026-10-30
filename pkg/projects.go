package pkg

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed resources/hackli.json
var hackliJsonContent string

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

func InitProject(path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	hackliJsonPath := filepath.Join(absolutePath, "hackli.json")
	if _, err := os.Stat(hackliJsonPath); os.IsNotExist(err) {
		file, err := os.Create(hackliJsonPath)
		if err != nil {
			return err
		}
		_, err = file.WriteString(hackliJsonContent)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		return fmt.Errorf("hackli.json already exists in project path: %s", absolutePath)
	}
	return nil
}
