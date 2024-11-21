package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"

	"ftp-scanner/config"
	"ftp-scanner/ftp"

	"bytes"
	"golang.org/x/net/html/charset"
	"io"
	"strings"
)

var dbConnect *sql.DB

func Connect(config config.Config) {
	db, err := sql.Open("postgres", connectionString(config))

	if err != nil {
		panic(err)
	}

	dbConnect = db
}

func Disconnect() {
	dbConnect.Close()
}

func AddServer(host string) int64 {
	serverId := 0
	sqlSelectStatement := `SELECT id FROM server where name = $1`

	var err = dbConnect.QueryRow(sqlSelectStatement, host).Scan(&serverId)
	if err == nil {
		return int64(serverId)
	}

	fmt.Println("Server not found trying to insert", host)

	sqlStatement := `INSERT INTO server(name) VALUES ($1) RETURNING id`
	err = dbConnect.QueryRow(sqlStatement, host).Scan(&serverId)
	if err == nil {
		return int64(serverId)
	}
	fmt.Println("Can't insert ", host)
	fmt.Println(err)

	return -1
}

func AddFiles(host string, wg *sync.WaitGroup, files []ftp.FileInfo) {
	defer wg.Done()

	serverId := AddServer(host)
	sqlStatement := `INSERT INTO file(full_name, fname, path, size, folder_id, server_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	sqlStatementFileAttr := `INSERT INTO file_attr(file_id, length) VALUES ($1, $2) RETURNING id`

	id := 0
	fileAttrId := 0
	for _, f := range files {
		fullName := f.Path + f.Name
		var err = dbConnect.QueryRow(sqlStatement, fullName, f.Name, f.Path, f.Size, *f.ParentId, serverId).Scan(&id)
		if err == nil {
			dbConnect.QueryRow(sqlStatementFileAttr, id, f.Size).Scan(&fileAttrId)
			continue
		}
		if strings.Contains(err.Error(), "fullname_unique") {
			continue
		}

		fmt.Println("1 Can't insert file with params (" + f.Name + "," + f.Path + "," + string(f.Size) + ")")
		fmt.Println(err)

		if strings.Contains(err.Error(), "invalid byte sequence for encoding") {
			err = dbConnect.QueryRow(sqlStatement, tidy(f.Path), tidy(f.Name), f.Path, f.Size).Scan(&id)
			if err == nil {
				continue
			}
			fmt.Println("2 Can't insert file with params (" + tidy(f.Path) + "," + tidy(f.Name) + "," + f.Path + "," + string(f.Size) + ")")
			fmt.Println(err)
		}
	}

}

func SaveFolder(host string, folder *ftp.FileInfo) {

	serverId := AddServer(host)

	id := 0
	sqlInsertStatement := `INSERT INTO folder(name, parent_id, server_id) VALUES ($1, $2, $3) RETURNING id`
	var err = dbConnect.QueryRow(sqlInsertStatement, folder.Path, folder.ParentId, serverId).Scan(&id)
	if err == nil {
		//fmt.Println("SaveFolder with id - ", id)
		folder.Id = uint64(id)
		return
	}

	if strings.Contains(err.Error(), "folder_name_unique") {
		sqlSelectStatement := `SELECT id FROM folder WHERE name = $1 and server_id = $2`
		err = dbConnect.QueryRow(sqlSelectStatement, folder.Path, serverId).Scan(&id)

		//fmt.Println("SaveFolder selected id - ", id)
		folder.Id = uint64(id)
		return
	}

	fmt.Printf("Can't insert or find ftp.folder with params (%s,%d, %d, %s)\n", folder.Path, folder.ParentId, folder.Id, serverId)
	fmt.Println("1 folder - ", folder.ParentId)
	fmt.Println("1 folder.Id - ", id)
	fmt.Println(err)
}

func AddFoldersToEmpty(host string, wg *sync.WaitGroup, files []ftp.FileInfo) {

	defer wg.Done()

	serverId := AddServer(host)
	sqlStatement := `INSERT INTO file(full_name, fname, path, size, folder_id, server_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	sqlStatementFileAttr := `INSERT INTO file_attr(file_id, length) VALUES ($1, $2) RETURNING id`

	id := 0
	fileAttrId := 0
	for _, f := range files {
		fullName := f.Path + f.Name
		var err = dbConnect.QueryRow(sqlStatement, fullName, f.Name, f.Path, f.Size, *f.ParentId, serverId).Scan(&id)
		if err == nil {
			dbConnect.QueryRow(sqlStatementFileAttr, id, f.Size).Scan(&fileAttrId)
			continue
		}
		if strings.Contains(err.Error(), "fullname_unique") {
			continue
		}

		fmt.Println("1 Can't insert file with params (" + f.Name + "," + f.Path + "," + string(f.Size) + ")")
		fmt.Println(err)

		if strings.Contains(err.Error(), "invalid byte sequence for encoding") {
			err = dbConnect.QueryRow(sqlStatement, tidy(f.Path), tidy(f.Name), f.Path, f.Size).Scan(&id)
			if err == nil {
				continue
			}
			fmt.Println("2 Can't insert file with params (" + tidy(f.Path) + "," + tidy(f.Name) + "," + f.Path + "," + string(f.Size) + ")")
			fmt.Println(err)
		}
	}
}

func UpdateFileAttrWithHash(config config.Config, fileId uint64, fileHash string) {
	db, err := sql.Open("postgres", connectionString(config))

	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `UPDATE file_attr SET hash=$1 WHERE file_id=$2`
	id := 0

	err = db.QueryRow(sqlStatement, fileHash, fileId).Scan(&id)
	if err != nil {
		panic(err)
	}
}

func QueryLargestFileNameWithouHash(config config.Config) (*uint64, string, uint64) {
	db, err := sql.Open("postgres", connectionString(config))

	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT f.id, full_name, max(size) size
                         FROM file f LEFT JOIN file_attr fa ON f.id=fa.file_id AND fa.hash IS NULL
                         GROUP BY f.id, full_name, path
                         ORDER BY  size DESC LIMIT 1`
	var id *uint64
	fullName := ""
	var size uint64
	//	size = 0

	var qId sql.NullInt64
	var qName sql.NullString
	var qSize sql.NullInt64

	err = db.QueryRow(sqlStatement).Scan(&qId, &qName, &qSize)
	if err != nil {
		fmt.Println(err.Error())
		if err == sql.ErrNoRows {
			// there were no rows, but otherwise no error occurred
			return nil, "", 0
		} else if strings.Contains(err.Error(), "destination pointer is nil") {
			// no rows skip
			return nil, "", 0
		} else {
			panic(err)

		}
	}

	fmt.Println("id ", qId, qName)
	if qId.Valid {
		t := uint64(qId.Int64)
		id = &t
	}
	if qName.Valid {
		fullName = qName.String
	}
	if qSize.Valid {
		size = uint64(qSize.Int64)
	}

	return id, fullName, size

}

func tidy(fileName string) string {
	strBytes := []byte(fileName)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel("windows-1251", byteReader)
	strBytes, _ = io.ReadAll(reader)
	return string(strBytes)
}

func connectionString(config config.Config) string {
	result := fmt.Sprintf(
		"host=%s "+
			"port=%s "+
			"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"sslmode=disable "+
			"search_path=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.Username,
		config.Database.Password,
		config.Database.DBName,
		config.Database.Schema)
	return result
}
