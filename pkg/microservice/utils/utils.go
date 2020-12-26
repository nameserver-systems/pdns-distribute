package utils

import (
	"crypto/sha256"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

func GenerateUUID() (string, error) {
	generatedUUID, err := uuid.NewUUID()
	if err != nil {
		return getHashedTime(), nil
	}

	uuidstring := generatedUUID.String()

	return uuidstring, err
}

func getHashedTime() string {
	now := time.Now().String()
	hash := sha256.Sum256([]byte(now))

	return string(hash[:])
}

func EnsurePathExist(path string) error {
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		return err
	}

	return nil
}
