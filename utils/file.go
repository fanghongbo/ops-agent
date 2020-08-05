package utils

import (
	"os"
	"path"
)

func FileIsExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

func IsFile(fp string) bool {
	f, err := os.Stat(fp)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

func IsDir(fp string) bool {
	f, err := os.Stat(fp)
	if err != nil {
		return false
	}
	return f.IsDir()
}

func WriteBytes(filePath string, b []byte) (int, error) {
	_ = os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = fw.Close()
	}()
	return fw.Write(b)
}

func WriteString(filePath string, s string) (int, error) {
	return WriteBytes(filePath, []byte(s))
}
