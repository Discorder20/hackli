package pkg

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Project struct {
	Name string `json:"name"`
}

//go:embed resources/hackli.json
var hackliJsonContent string

func VerifyStructure(path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	//check if the path exists
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", absolutePath)
	}
	// check if it includes hackli.json
	hackliJsonPath := filepath.Join(absolutePath, "hackli.json")
	if _, err := os.Stat(hackliJsonPath); os.IsNotExist(err) {
		return fmt.Errorf("hackli.json not found in project path: %s", absolutePath)
	}
	return nil
}

func GetProject(path string) (*Project, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	verifyErr := VerifyStructure(absolutePath)
	if verifyErr != nil {
		return nil, verifyErr
	}
	jsonContent, err := os.ReadFile(filepath.Join(absolutePath, "hackli.json"))
	if err != nil {
		return nil, err
	}

	var config Project
	if err := json.Unmarshal(jsonContent, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func OpenProject(path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	verifyErr := VerifyStructure(absolutePath)
	if verifyErr != nil {
		return verifyErr
	}
	project, err := GetProject(absolutePath)
	if err != nil {
		return err
	}
	return OpenEditor(absolutePath, project)
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
