package report

import (
	"os"
)

func New(outPath string, data []byte) error {

	fd, err := os.Create(outPath)
	defer fd.Close()
	if err != nil {
		return err
	}
	_, err = fd.Write(data)
	if err != nil {
		return err
	}
	return nil
}
