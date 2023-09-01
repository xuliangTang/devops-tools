package helpers

import (
	"io"
	"log"
	"os"
)

func MustLoadFile(path string) []byte {
	b, err := LoadFile(path)
	if err != nil {
		log.Fatalln("文件不存在:", path)
	}

	return b
}

func LoadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	return b, err
}

func MustFileExists(file string) {
	if !FileExists(file) {
		panic("文件:" + file + "不存在")
	}
}

func FileExists(file string) bool {
	_, err := os.Stat(file)

	if os.IsNotExist(err) {
		return false
	}

	return true
}
