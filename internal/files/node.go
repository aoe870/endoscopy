package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/klauspost/compress/zip"
	"github.com/mholt/archiver/v4"
	"github.com/pkg/errors"
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
	ff, _ := os.Open(path)
	fsys, err := FileSystem(context.Background(), path, ff)
	if err != nil {
		return nil, err
	}

	err = fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
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

		info, err := d.Info()
		// 软链接类型不处理
		if info.Mode().Type() == fs.ModeSymlink {
			return nil
		}
		buff := bytes.NewBuffer(make([]byte, 0, info.Size()))
		fn, err := fsys.Open(p)
		if err != nil {
			return nil
		}
		buff.ReadFrom(fn)
		fn.Close()
		format, _, err := archiver.Identify(p, fn)
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
		fmt.Println("Walking:", path, "Dir?", d.IsDir())

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

func FileSystem(ctx context.Context, root string, file *os.File) (fs.FS, error) {

	info, err := file.Stat()
	format, _, err := archiver.Identify(filepath.Base(root), file)
	if err != nil && !errors.Is(err, fmt.Errorf("no formats matched")) {
		return nil, err
	}

	if format != nil {
		switch ff := format.(type) {
		case archiver.Zip:
			// zip.Reader is more performant than ArchiveFS, because zip.Reader caches content information
			// and zip.Reader can open several content files concurrently because of io.ReaderAt requirement
			// while ArchiveFS can't.
			// zip.Reader doesn't suffer from issue #330 and #310 according to local test (but they should be fixed anyway)

			// open the file anew, as our original handle will be closed when we return
			file, err := os.Open(root)
			if err != nil {
				return nil, err
			}
			return zip.NewReader(file, info.Size())
		case archiver.Archival:
			// TODO: we only really need Extractor and Decompressor here, not the combined interfaces...
			return archiver.ArchiveFS{Path: root, Format: ff, Context: ctx}, nil
		case archiver.Compression:
			return archiver.FileFS{Path: root, Compression: ff}, nil
		}
	}

	// otherwise consider it an ordinary file; make a file system with it as its only file
	return archiver.FileFS{Path: root}, nil
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
