package ftp

import (
	"ftp-scanner/config"
	//	"fmt"
	"github.com/secsy/goftp"
	"os"
	"strings"
)

var client *goftp.Client

func Connect(config config.Config) {
	var ftpConfig = goftp.Config{
		User:     config.Ftp.Username,
		Password: config.Ftp.Password,
	}
	c, err := goftp.DialConfig(ftpConfig, config.Ftp.Host)
	if err != nil {
		panic(err)
	}

	client = c
}

func Disconnect() {
	client.Close()
}

func ReadDir(baseFolder FileInfo) (*[]FileInfo, *[]FileInfo, *[]FileInfo) {
	var filesResult = new([]FileInfo)
	var foldersResult = new([]FileInfo)
	var emptyFolder = new([]FileInfo)

	var dir = baseFolder.Path
	files, err := client.ReadDir(dir)

	if err != nil {
		return filesResult, foldersResult, emptyFolder
	}
	if len(files) == 0 {
		var finfo = FileInfo{Name: "", Path: dir, IsDir: true, Size: 0, ParentId: &baseFolder.Id, FileInfo: &baseFolder}
		*emptyFolder = append(*emptyFolder, finfo)
		return filesResult, foldersResult, emptyFolder
	}

	for _, f := range files {
		if f.IsDir() {
			fullDir := dir + f.Name() + "/"
			if dir == "/" {
				fullDir = "/" + f.Name() + "/"
			}
			var finfo = FileInfo{Name: f.Name(), Path: fullDir, IsDir: f.IsDir(), Size: f.Size(), ParentId: &baseFolder.Id, FileInfo: &baseFolder}
			*foldersResult = append(*foldersResult, finfo)
		} else {
			var finfo = FileInfo{Name: f.Name(), Path: baseFolder.Path, IsDir: f.IsDir(), Size: f.Size(), ParentId: &baseFolder.Id, FileInfo: &baseFolder}
			*filesResult = append(*filesResult, finfo)
		}

	}

	return filesResult, foldersResult, emptyFolder
}

func DownloadAsBuffer() []int {
	var result = []int{}
	return result
}

func DownloadAsFile(file string) {
	temp, err := os.Create(".tmp")
	if err != nil {
		panic(err)
	}

	err = client.Retrieve(file, temp)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			//skip
		}
		panic(err)
	}

}
