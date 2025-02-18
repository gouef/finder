package finder

import "os"

type Info struct {
	FileInfo os.FileInfo
	Path     string
	Ext      string
	Name     string
}
