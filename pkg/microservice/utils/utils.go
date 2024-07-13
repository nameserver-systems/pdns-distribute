package utils

import (
	"os"
	"strconv"
	"strings"
)

func TrimAndLowerString(in string) string {
	output := trimSpace(in)
	output = strings.ToLower(output)

	return output
}

func trimSpace(text string) string {
	output := strings.TrimSpace(text)
	output = strings.ReplaceAll(output, " ", "")

	return output
}

func ConvertStringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)

	return i, err
}

func EnsurePathExist(path string) error {
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		return err
	}

	return nil
}
