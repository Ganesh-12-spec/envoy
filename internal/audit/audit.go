package audit

import (
	"fmt"
	"os"
	"time"

	"github.com/Ganesh-12-spec/envoy/internal/paths"
)

func Log(action string, target string) error {
	logPath := paths.Audit()

	if err := os.MkdirAll(paths.EnvoyDir(), 0700); err != nil {
		return fmt.Errorf("creating Envoy directory: %w", err)
	}

	file, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0600,
	)
	if err != nil {
		return fmt.Errorf("opening audit log: %w", err)
	}
	defer file.Close()

	timestamp := time.Now().UTC().Format(time.RFC3339)

	entry := timestamp + " " + action

	if target != "" {
		entry += " " + target
	}

	entry += "\n"

	if _, err := file.WriteString(entry); err != nil {
		return fmt.Errorf("writing audit log: %w", err)
	}

	return nil
}
