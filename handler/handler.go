package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "inner error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Errorf("fail to get file,err:%s", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: "/tmp/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		create, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("fail to create file,err:%s", err.Error())
			return
		}
		defer create.Close()
		fileMeta.FileSize, err = io.Copy(create, file)
		if err != nil {
			fmt.Printf("faile to copy file,err:%s", err.Error())
			return
		}
		create.Seek(0, 0)
		all, err := ioutil.ReadAll(create)
		fileMeta.FileSha1 = util.Sha1(all)
		meta.UpdateFileMetas(fileMeta)
		fmt.Println(fileMeta)
		fmt.Println(meta.FileMetas)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}

}

//UploadSuccessHandler : 上传完成
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload finished!")
}

//GetFileMetaHandler 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	s := r.Form["filehash"][0]
	metas := meta.GetFileMetas(s)
	data, err := json.Marshal(metas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	s := r.Form["filehash"][0]
	metas := meta.GetFileMetas(s)
	open, err := os.Open(metas.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer open.Close()
	data, err := ioutil.ReadAll(open)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+metas.FileName+"\"")
	w.Write(data)
}

func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileSha1 := r.Form["filehash"][0]
	fileName := r.Form["filename"][0]
	opType := r.Form["op"][0]

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metas := meta.GetFileMetas(fileSha1)
	metas.FileName = fileName
	meta.UpdateFileMetas(metas)
	marshal, err := json.Marshal(metas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form["filehash"][0]

	metas := meta.GetFileMetas(fileSha1)
	os.Remove(metas.Location)

	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}
