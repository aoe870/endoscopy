package files

import (
	"fmt"
	"github.com/mholt/archiver/v4"
	"io/fs"
	"os"
	"path/filepath"
)

type FileType int

const (
	File     FileType = 1
	Ctalogue FileType = 2
	Archived FileType = 3
)

type FileMetadata struct {
	Name     string
	FileType FileType
	Table    []*Node
}

func New(path string) (FileMetadata, error) {

	var file FileMetadata

	_, err := os.Stat(path)
	if err != nil {
		return file, err
	}

	file.readFile(path)

	return file, nil
}

func (metadata FileMetadata) readFile(path string) {
	f, _ := os.Stat(path)
	if f.IsDir() {
		err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			Node := Node{
				Name: d.Name(),
				Path: p,
			}
			metadata.Table = append(metadata.Table, &Node)
			if d.IsDir() {
				Node.FileType = Ctalogue
				return nil
			}
			Node.FileType = File

			// 判断是否为压缩文件
			input, _ := os.ReadFile(p)
			format, input, err := archiver.Identify(p, input)
			if err != nil {
				return err
			}

			if ex, ok := format.(archiver.Extractor); ok {
				// ... proceed to extract
			}

			// or maybe it's compressed and you want to decompress it?
			if decom, ok := format.(archiver.Decompressor); ok {
				rc, err := decom.OpenReader(unknownFile)
				if err != nil {
					return err
				}
				defer rc.Close()

				// read from rc to get decompressed data
			}

			return nil
		})

		if err != nil {
			fmt.Println(err)
		}

	} else {
		// 读取文件

	}
}
