package files

import (
	"bytes"
	"github.com/klauspost/compress/zip"
	"github.com/mholt/archiver/v4"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
		node := Node{
			Name: d.Name(),
			Path: p,
		}
		metadata.Table = append(metadata.Table, &node)
		if d.IsDir() {
			node.FileType = Ctalogue
			return nil
		}
		node.FileType = File
		// 判断是否为压缩文件
		f, _ := os.Open(p)
		fn, _ := f.Stat()
		defer f.Close()

		//zip单独处理
		if filepath.Ext(p) == ".jar" ||
			filepath.Ext(p) == ".war" {

			//判断文件大小
			if fn.Size() > 1*1024*1024*1024 {
				return nil
			}

			//提取jar文件文件名，不包含jar后缀
			zr, err := zip.NewReader(f, fn.Size())
			if err != nil {
				return nil
			}
			for _, zf := range zr.File {
				if zf.Mode().IsDir() {
					continue
				}
				fr, err := zf.Open()
				defer fr.Close()
				if err != nil {
					return nil
				}
				buff := bytes.NewBuffer(make([]byte, 0, zf.FileInfo().Size()))
				buff.ReadFrom(fr)
				if strings.Contains(zf.Name, ".jar") {
					continue
				}
				metadata.Table = append(metadata.Table, &Node{
					Name:     zf.Name,
					Path:     filepath.Base(p) + ":" + zf.Name,
					FileType: File,
					Data: &FileData{
						Data: buff.Bytes(),
					},
				})
			}
			return nil
		}

		_, _, err = archiver.Identify(p, f)
		if err != nil {
			if fn.Size() > 1*1024*1024 {
				return nil
			}
			buff := bytes.NewBuffer(make([]byte, 0, fn.Size()))
			buff.ReadFrom(f)
			node.Data = &FileData{
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
			node.FileType = Archived
		}

		return nil
	})

	if err != nil {
		return
	}

	return
}
