package report

import (
	"os"
)

func New(outPath string, data []byte) {

	fd, err := os.Create(outPath)
	defer fd.Close()
	if err != nil {
		return
	}
	_, err = fd.Write(data)
	if err != nil {
		return
	}
	return
}
