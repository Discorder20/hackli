package pkg

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func OpenEditor(path string, project *Project) error {
	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("enable terminal input: %w", err)
	}
	defer term.Restore(fd, state)

	fmt.Print("\033[H\033[2JEditor Opened\r\nPress q to quit.\r\n")
	fmt.Printf("Project: %s\r\n", project.Name)

	typingMode := false
	var key [1]byte
	for {
		n, err := os.Stdin.Read(key[:])
		if err != nil {
			return fmt.Errorf("read terminal input: %w", err)
		}
		if n > 0 && key[0] == 'q' && !typingMode {
			return nil
		}
	}
}
