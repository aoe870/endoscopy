package files

import (
	"bufio"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/mholt/archiver/v4"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type FileMetadata struct {
	Name     string
	FileType FileType
	Table    []*Node
}

func New(path string) (*FileMetadata, error) {

	tempDir, _ := ioutil.TempDir("", fmt.Sprintf("sca-%x", md5.Sum([]byte(path+strconv.FormatInt(time.Now().Unix(), 10)))))
	f, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	locationPath := strings.ReplaceAll(path, ":", "")

	var fileNodes []*Node
	if f.IsDir() {
		locationPath = path
		locationPath = filepath.Join(tempDir, strings.ReplaceAll(locationPath, ":", ""))
		err := copyFolder(path, locationPath)
		if err != nil {
			return nil, err
		}

		fileNode, _ := walkPath(path, locationPath, tempDir)
		fileNodes = append(fileNodes, fileNode...)
	} else {

		node := Node{
			Name:         f.Name(),
			Path:         path,
			LocationPath: locationPath,
			FileType:     File,
		}

		isIdentity, Type := getIdentityFile(path)
		if isIdentity {
			locationPath, err = decompressor(path, filepath.Join(tempDir, locationPath))
			node.LocationPath = locationPath
			if err != nil {
				return nil, err
			}
			node.FileType = Type
			fileNode, _ := walkPath(path, locationPath, tempDir)
			fileNodes = append(fileNodes, fileNode...)

		} else {
			data, _ := os.ReadFile(path)
			node.Data = &FileData{
				Data: data,
			}

			node.LocationPath = filepath.Join(tempDir, locationPath)
			node.SaveFile(node.LocationPath)
		}

		fileNodes = append(fileNodes, &node)
	}

	fileMetadata := &FileMetadata{
		Name:  path,
		Table: fileNodes,
	}

	return fileMetadata, nil
}

func copyFolder(source, destination string) error {
	// 创建目标文件夹
	err := os.MkdirAll(destination, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination folder: %v", err)
	}

	// 遍历源文件夹
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %q: %v", path, err)
		}

		// 生成目标文件/文件夹路径
		destPath := filepath.Join(destination, path[len(source):])

		if info.IsDir() {
			// 如果是文件夹，则在目标路径中创建一个对应的文件夹
			err = os.MkdirAll(destPath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory at %q: %v", destPath, err)
			}
		} else {
			// 如果是文件，则复制文件到目标路径
			err = copyFile(path, destPath)
			if err != nil {
				return fmt.Errorf("failed to copy file from %q to %q: %v", path, destPath, err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy folder: %v", err)
	}

	return nil
}

func copyFile(source, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy data: %v", err)
	}

	return nil
}

// 解析文件/文件夹
func walkPath(filePath, locationPath, temp string) ([]*Node, error) {

	fsys, err := archiver.FileSystem(nil, locationPath)
	if err != nil {
		return nil, err
	}

	var fileNode []*Node
	err = fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {

		if p == "." {
			return nil
		}

		info, _ := d.Info()
		if info.Mode().Type() == fs.ModeSymlink {
			return nil
		}
		node := Node{
			Name:         d.Name(),
			Path:         filepath.Join(filePath, p),
			LocationPath: filepath.Join(locationPath, p),
			IsFile:       !d.IsDir(),
		}

		if d.IsDir() {
			node.FileType = Catalog
			fileNode = append(fileNode, &node)
			return nil
		}
		isIdentity, Type := getIdentityFile(p)
		if isIdentity {
			archiverPath := filepath.Join(filePath, p)
			path := filepath.Join(locationPath, p)
			location, _ := decompressor(path, path)
			archive, _ := walkPath(archiverPath, location, temp)
			node.FileType = Type
			fileNode = append(fileNode, archive...)
		} else {
			data, err := os.ReadFile(node.LocationPath)
			if err != nil {
				return nil
			}
			node.Data = &FileData{
				Data: data,
			}
		}
		fileNode = append(fileNode, &node)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileNode, nil
}

func decompressor(source, dst string) (string, error) {

	p := filepath.Join(filepath.Dir(dst), strings.ReplaceAll(filepath.Base(dst), ".", ""))
	if isIdentityFile(source, JavaArchived, Apk) {
		err := unzipSource(source, p)
		if err != nil {
			return "", err
		}
		return p, nil
	}

	f, err := os.Open(source)
	if err != nil {
		return "", err
	}
	format, input, err := archiver.Identify(source, f)
	if err != nil {
		return "", err
	}
	// you can now type-assert format to whatever you need;
	// be sure to use returned stream to re-read consumed bytes during Identify()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 定义并发任务数量
	concurrency := 5

	// 创建一个wait group来同步任务
	var wg sync.WaitGroup

	// 创建一个channel来控制并发任务数量
	semaphore := make(chan struct{}, concurrency)

	//var wg sync.WaitGroup
	handler := func(ctx context.Context, f archiver.File) error {

		tempPath := filepath.Join(p, f.NameInArchive)
		if !utf8.ValidString(f.NameInArchive) {
			gbkName := mahonia.NewDecoder("gbk").ConvertString(f.NameInArchive)
			_, cdata, _ := mahonia.NewDecoder("utf-8").Translate([]byte(gbkName), true)
			tempPath = filepath.Join(p, string(cdata))
		}
		if f.IsDir() {
			err := os.MkdirAll(tempPath, os.ModePerm)
			if err != nil {
				return nil
			}
		} else {
			path := filepath.Dir(tempPath)
			parentF, err := os.Stat(path)
			if parentF == nil {
				err = os.MkdirAll(path, os.ModePerm)
				if err != nil {
					return nil
				}
			}

			f, err := f.Open()
			if err != nil {

				return nil
			}
			defer f.Close()
			r := bufio.NewReader(f)
			data, err := ioutil.ReadAll(r)
			if err != nil {
				return nil
			}

			// 在执行任务之前获取一个信号量
			semaphore <- struct{}{}
			// 增加wait group的计数器
			wg.Add(1)
			go func(data []byte, path string) {
				// 在任务完成之后释放信号量
				defer func() {
					<-semaphore
				}()

				defer wg.Done()
				targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				if err != nil {
					return
				}
				defer targetFile.Close()

				_, err = targetFile.Write(data)
				if err != nil {
					return
				}
			}(data, tempPath)
		}

		// do something with the file
		return nil
	}
	// want to extract something?
	if ex, ok := format.(archiver.Extractor); ok {
		// ... proceed to extract
		err = ex.Extract(ctx, input, nil, handler)
		if err != nil {
			return "", err
		}

	}
	wg.Wait()

	return p, nil
}

func (node *Node) SaveFile(path string) bool {
	if node.Data == nil {
		return false
	}

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return false
	}

	fd, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return false
	}
	defer fd.Close()

	if err := os.WriteFile(path, node.Data.Data, 0666); err != nil {
		return false
	}
	node.LocationPath = path
	return true
}

func unzipSource(source, destination string) error {
	// 1. Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		if f.Name == "/" {
			continue
		}
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}
