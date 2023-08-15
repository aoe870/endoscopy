package files

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/klauspost/compress/zip"
	"github.com/mholt/archiver/v4"
)

type Node struct {
	Name     string
	Path     string
	FileType FileType
	Data     *FileData
}

type FileData struct {
	Data []byte
}

func (node *Node) readNode() {
	if node.FileType == Ctalogue {
		return
	}
}

func readArchives(prePath, path string) ([]*Node, error) {

	var fileList []*Node
	system, err := archiver.FileSystem(context.Background(), path)
	if err != nil {
		return nil, err
	}
	err = fs.WalkDir(system, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if p == ".git" {
			return fs.SkipDir
		}

		if p == ".idea" {
			return fs.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		if p == "." {
			return nil
		}

		info, err := d.Info()
		// 软链接类型不处理
		if info.Mode().Type() == fs.ModeSymlink {
			return nil
		}
		buff := bytes.NewBuffer(make([]byte, 0, info.Size()))
		fn, err := system.Open(p)
		if err != nil {
			return nil
		}
		buff.ReadFrom(fn)
		fn.Close()

		//zip单独处理
		if filepath.Ext(p) == ".jar" ||
			filepath.Ext(p) == ".war" {

			//判断文件大小
			if info.Size() > 1*1024*1024*1024 {
				return nil
			}

			//提取jar文件文件名，不包含jar后缀
			reader := bytes.NewReader(buff.Bytes())
			buf := make([]byte, info.Size())
			if _, err := reader.Read(buf); err != nil {
				return nil
			}
			zr, err := zip.NewReader(reader, info.Size())
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
				fileList = append(fileList, &Node{
					Name:     zf.Name,
					FileType: File,
					Path:     p + ":" + zf.FileInfo().Name(),
					Data:     &FileData{buff.Bytes()},
				})
			}
			return nil
		}

		newFile := NewFile()
		_, err = newFile.Write(buff.Bytes())
		format, _, err := archiver.Identify(p, newFile)
		newFile.Close()
		if format != nil {
			path, _ := ioutil.TempDir("", "endoscopy-"+strings.ReplaceAll(p, "/", "-"))
			path = filepath.Join(path, d.Name())
			_, err := createFile(path, buff.Bytes())

			if err == nil {
				table, err := readArchives(p, path)
				if err == nil {
					fileList = append(fileList, table...)
				}
			}
			os.RemoveAll(filepath.Dir(path))
		}
		node := &Node{
			Name:     d.Name(),
			FileType: File,
			Path:     p,
			Data:     &FileData{buff.Bytes()},
		}
		if prePath != "" {
			node.Path = prePath + ":" + node.Path
		}
		fileList = append(fileList, node)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return fileList, nil
}

// 读取文件
func decompressor(filename string) (string, error) {

	temp := createTempDir()
	fsys, err := archiver.FileSystem(nil, filename)
	if err != nil {
		return temp, nil
	}

	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == ".git" ||
			path == ".idea" ||
			path == ".svn" {
			return fs.SkipDir
		}

		tempPath := filepath.Join(temp, path)
		if d.IsDir() {
			err = os.MkdirAll(tempPath, os.ModePerm)
			if err != nil {
				return nil
			}
		} else {
			fr, err := fsys.Open(path)
			defer fr.Close()
			if err != nil {
				return nil
			}
			// 创建要写出的文件对应的 Write
			fw, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
			defer fw.Close()
			if err != nil {
				return nil
			}

			_, err = io.Copy(fw, fr)
			if err != nil {
				return nil
			}
		}

		return nil
	})
	if err != nil {
		return temp, err
	}
	return temp, nil
}

func createFile(path string, data []byte) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	f.Write(data)
	defer f.Close()
	return f, nil
}

func createTempDir() string {
	name := "endoscopy" + time.Now().Format("20060102150405")
	path, err := ioutil.TempDir("", name)
	if err != nil {
		return createTempDir()
	}
	return path
}
