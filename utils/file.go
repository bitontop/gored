package utils

import (
	"fmt"
	"os"
)

func SaveFile(path, str string) error {
	f, err := os.OpenFile(path,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", str)); err != nil {
		return err
	}
	return nil
}
