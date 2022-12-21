package cfg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	separator = "="
	envPath   = ".env"
)

func Load() (err error) {
	file, err := os.Open(envPath)
	if err != nil {
		return fmt.Errorf("error loading %s: %w", envPath, err)
	}

	defer func() {
		if e := file.Close(); e != nil {
			err = fmt.Errorf("error closing %s: %w", envPath, e)
		}
	}()

	buf := bufio.NewScanner(file)
	buf.Split(bufio.ScanLines)

	for buf.Scan() {
		if keyVal := strings.Split(buf.Text(), separator); keyVal != nil {
			if err := os.Setenv(keyVal[0], keyVal[1]); err != nil {
				return fmt.Errorf("error setting environment variable: %w", err)
			}
		}
	}

	return nil
}
