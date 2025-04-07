package files

import (
	"path/filepath"
	"strings"
)

var skipDir = []string{
	"__MACOSX",
	".git",
	".svn",
}

type parseEntry struct {
	Name string
	Type FileType
}

var identityFiles = []parseEntry{
	{
		Name: ".zip",
		Type: Archived,
	},
	{
		Name: ".tar",
		Type: Archived,
	},
	{
		Name: ".rar",
		Type: Archived,
	},
	{
		Name: ".7z",
		Type: Archived,
	},
	{
		Name: ".br",
		Type: Archived,
	},
	{
		Name: ".bz2",
		Type: Archived,
	},
	{
		Name: ".gz",
		Type: Archived,
	},
	{
		Name: ".lz4",
		Type: Archived,
	},
	{
		Name: ".sz",
		Type: Archived,
	}, {
		Name: ".xz",
		Type: Archived,
	}, {
		Name: ".zz",
		Type: Archived,
	},
	{
		Name: ".zst",
		Type: Archived,
	},
	{
		Name: ".jar",
		Type: JavaArchived,
	},
	{
		Name: ".war",
		Type: JavaArchived,
	},
	{
		Name: ".apk",
		Type: Apk,
	},
}

type FileType int32

const (
	Unknown      FileType = 1  // 未知格式
	File         FileType = 2  // 文件
	Catalog      FileType = 3  // 目录
	Binary       FileType = 4  // 二进制
	JavaArchived FileType = 5  // Jar/War
	Archived     FileType = 6  // 存档文件, 压缩文件
	GoBinary     FileType = 8  // go二进制
	Apk          FileType = 9  //安卓包
	LinkTarget   FileType = 10 //链接
	IPA          FileType = 11 //IOS包
	ARR          FileType = 12 //android arr包
)

func getIdentityFile(path string) (bool, FileType) {
	for _, fileType := range identityFiles {
		base := filepath.Base(path)
		if strings.HasSuffix(base, fileType.Name) {
			return true, fileType.Type
		}
	}
	return false, Unknown
}

func isIdentityFile(path string, fileType ...FileType) bool {
	allType := make(map[FileType]struct{})
	for _, fileType := range fileType {
		allType[fileType] = struct{}{}
	}

	for _, fileType := range identityFiles {
		if _, ok := allType[fileType.Type]; !ok {
			continue
		}
		base := filepath.Base(path)
		if strings.HasSuffix(base, fileType.Name) {
			return true
		}
	}

	return false
}
