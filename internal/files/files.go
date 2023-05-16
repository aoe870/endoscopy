package files

import (
	"bytes"
	"github.com/mholt/archiver/v4"
	"io/fs"
	"os"
	"path/filepath"
)

type FileType int

const (
	File       FileType = 1
	Ctalogue   FileType = 2
	Archived   FileType = 3
	BinaryFile FileType = 4
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

	file.readFile("", path)

	return file, nil
}

func (metadata *FileMetadata) readFile(prePath, path string) {

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
		f, _ := os.Open(p)
		fn, _ := f.Stat()
		defer f.Close()
		_, _, err = archiver.Identify(p, f)
		if err != nil {
			buff := bytes.NewBuffer(make([]byte, 0, fn.Size()))
			buff.ReadFrom(f)
			Node.Data = &FileData{
				buff.Bytes(),
			}
		} else {
			//判断文件大小
			if fn.Size() > 1*1024*1024*1024 {
				tempPath, _ := decompressor(p)
				metadata.readFile(p, tempPath)
				defer os.Remove(tempPath)
				return nil
			}
			//读取压缩文件
			table, err := readArchives(prePath, p)
			if err != nil {
				return err
			}
			metadata.Table = append(metadata.Table, table...)
			Node.FileType = Archived
		}

		return nil
	})

	if err != nil {
		return
	}

	return
}
