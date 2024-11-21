package db

type File struct {
	Id       uint64
	FullName string
	FileName string
	Path     string
	Size     int64
	FolderId uint64
	ServerId uint64
}

type FileAttr struct {
	Id     uint64
	FileId uint64
	Hash   string
	Length int64
}

type Folder struct {
	Id       uint64
	Name     string
	ParentId uint64
	ServerId uint64
}

type Server struct {
	Id   uint64
	Name uint64
}
