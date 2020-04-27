package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

func uploaderHandler(w http.ResponseWriter, req *http.Request) {

	//we first extract the form fields from the req object
	userId := req.FormValue("userid")

	//file contains the file in form of multi part bytes
	//header contains the info about file like filename and other meta data
	//err is the error
	file, header, err := req.FormFile("avatarFile")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//we read the file using ioutil package
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//we create a new filename using userId and existing filename from the header.Filename of file
	filename := path.Join("avatars", userId+path.Ext(header.Filename))
	//ioutil.WriteFile creates a new file in the avatars folder
	//0777 gives us all file permissions
	err = ioutil.WriteFile(filename, data, 0777)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	io.WriteString(w, "Successful")
}
