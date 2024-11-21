package ftp

type FileInfo struct {
	Id       uint64
	Name     string
	Path     string
	IsDir    bool
	Size     int64
	ParentId *uint64
	FileInfo *FileInfo
}

type FileRecord map[string]FileInfo
