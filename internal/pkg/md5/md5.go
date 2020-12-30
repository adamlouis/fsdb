package md5

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// GetMD5 returns the md5 sum of the file at the provided path
func GetMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sum := md5.New()
	if _, err := io.Copy(sum, file); err != nil {
		return "", err
	}

	bytes := sum.Sum(nil)[:16]
	return hex.EncodeToString(bytes), nil
}
