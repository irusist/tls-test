package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/upload", fileUploadHandler)

	fmt.Println("Server started at http://localhost:8899")
	http.ListenAndServe(":8899", nil)
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// 返回一个简单的文件上传表单
		fmt.Fprintf(w, `
            <html>
            <head>
              <title>Upload file</title>
            </head>
            <body>
              <h1>Upload file</h1>
              <form enctype="multipart/form-data" action="/upload" method="post">
                <input type="file" name="file" />
                <input type="submit" value="Upload" />
              </form>
            </body>
            </html>`)
		return
	} else if r.Method == http.MethodPost {
		// 解析请求的multipart/form-data类型
		if err := r.ParseMultipartForm(10 << 20); err != nil { // Max upload size ~10MB
			http.Error(w, "The uploaded file is too big. Please choose an file less than 10MB in size", http.StatusBadRequest)
			return
		}

		// 从表单中获取file字段
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving the File", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		// 读取文件内容
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading the file", http.StatusInternalServerError)
			return
		}

		// 随意处理文件（如下所示写入磁盘）
		tempFile, err := ioutil.TempFile("upload", "upload-*.tmp")
		if err != nil {
			log.Println("Error creating the temp file", err)
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		// 将文件写入服务器
		tempFile.Write(fileBytes)

		// 返回响应
		fmt.Fprintf(w, "Successfully Uploaded File\n")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
