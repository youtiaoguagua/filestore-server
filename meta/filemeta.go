package meta

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var FileMetas map[string]FileMeta

func init() {
	FileMetas = make(map[string]FileMeta)
}

//UpdateFileMetas :新增/更新元信息
func UpdateFileMetas(fmeta FileMeta) {
	FileMetas[fmeta.FileSha1] = fmeta
}

//GetFileMetas 获取元信息
func GetFileMetas(hash string) FileMeta {
	return FileMetas[hash]
}

//func GetLastFileMetas(count int) []FileMeta {
//	var fileMetaList []FileMeta
//	for _,v := range FileMetas {
//		fileMetaList = append(fileMetaList, v)
//	}
//	sort.Sort((fileMetaList))
//	return fileMetaList
//}

// RemoveFileMeta : 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(FileMetas, fileSha1)
}