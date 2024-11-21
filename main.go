package main

import (
	"flag"
	"fmt"
	"ftp-scanner/config"
	"ftp-scanner/db"
	"ftp-scanner/ftp"
	"ftp-scanner/hash"
	"github.com/secsy/goftp"
	"os"
	"sync"
	"time"
)

var mset ftp.FileRecord

func main() {
	startDir := flag.String("start-dir", "/", "directory on ftp to start scanning files")
	flag.Parse()
	fmt.Println(*startDir)

	var wg = new(sync.WaitGroup)

	ev := config.Read("config.yml")

	value := db.AddServer(ev)
	fmt.Println("ServerId = ", value)
	//if true {
	//	return
	//}

	mset := ftp.FileRecord{}
	ftpClient := createFtpClient(ev)

	var startPath = *startDir
	var lastSymbol = (*startDir)[len(*startDir)-1:]
	if lastSymbol == "/" {
		startPath = (*startDir)[0 : len(*startDir)-1]
		if len(startPath) == 0 {
			startPath = "/"
		}
	}

	var fileInfo = ftp.FileInfo{Name: *startDir, Path: startPath, IsDir: true, Size: 0, ParentId: nil}

	db.SaveFolder(ev, &fileInfo)
	mset[*startDir] = fileInfo

	//wg.Add(1)
	//go hashFiles(ev, wg)

	dirs := 0
	nums := 0
	for len(mset) > 0 {
		for k, v := range mset {
			var files, folders, emptyFolders = ftp.ReadDir(ftpClient, v)
			for _, folder := range *folders {
				db.SaveFolder(ev, &folder)
				var key = folder.Path + folder.Name
				mset[key] = folder
			}
			dirs++

			nums = nums + len(*files)
			fmt.Print("\r Total dirs \t", dirs, "\t|\t Total files ", nums, "             ")
			delete(mset, k)
			wg.Add(1)
			//                fmt.Println("Scanned ", k)
			go db.AddFiles(ev, wg, *files)
			wg.Add(1)
			go db.AddFoldersToEmpty(ev, wg, *emptyFolders)
		}
		fmt.Println("\n====================  Map size =", len(mset))
	}
	fmt.Println(mset)
	defer func(ftpClient *goftp.Client) {
		_ = ftpClient.Close()
	}(ftpClient)
	wg.Wait()
}

func createFtpClient(conf config.Config) *goftp.Client {
	var ftpConfig = goftp.Config{
		User:     conf.Ftp.Username,
		Password: conf.Ftp.Password,
	}
	client, err := goftp.DialConfig(ftpConfig, conf.Ftp.Host)
	if err != nil {
		panic(err)
	}

	return client
}

func hashFiles(conf config.Config, wg *sync.WaitGroup) {
	defer wg.Done()
	ftpClient := createFtpClient(conf)

	var id *uint64
	var fullName string
	var size uint64

	id, fullName, size = db.QueryLargestFileNameWithouHash(conf)
	retry := 0
	for retry < 10 {
		fmt.Println("\nFullName", fullName)
		if id == nil {
			retry++
			time.Sleep(3 * time.Second)
		} else {
			ftp.DownloadAsFile(ftpClient, fullName)
			fi, err := os.Stat(".tmp")
			var fileHash = "wrong"
			if err == nil && size == uint64(fi.Size()) {

				// get the size

				fileHash = hash.CalcHash(".tmp")
				//			db.InsertFileAttr(conf, *id, h, size)
			}
			db.UpdateFileAttrWithHash(conf, *id, fileHash)
			os.Remove(".tmp")

		}
		id, fullName, _ = db.QueryLargestFileNameWithouHash(conf)
	}
}
