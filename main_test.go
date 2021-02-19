package main

import(
	"testing"
	"os"
	"bytes"
	"mime/multipart"
	"path/filepath"
	"net/http/httptest"
	"net/http"
	"io"
	"github.com/stretchr/testify/assert"
)

func TestUploadTest(t *testing.T){
	assert := assert.New(t)
	path := "./icon.png"
	
	file, _ := os.Open(path)
	defer file.Close()
	
	os.RemoveAll("./uploads")

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	multi, err := writer.CreateFormFile("upload_file",filepath.Base(path))
	assert.NoError(err)
	
	io.Copy(multi,file)
	writer.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST","/uploads",buf)
	req.Header.Set("Content-type",writer.FormDataContentType())
	uploadsHandler(res,req)
	assert.Equal(http.StatusOK, res.Code)

	uploadFilePath := "./uploads/"+ filepath.Base(path)
	_, err1 := os.Stat(uploadFilePath)
	assert.NoError(err1)
	
	uploadFile, _ := os.Open(uploadFilePath)
	originFile, _ := os.Open(path)
	defer uploadFile.Close()
	defer originFile.Close()

	uploadData := []byte{}
	originData := []byte{}
	uploadFile.Read(uploadData)
	originFile.Read(originData)

	assert.Equal(originData,uploadData)
}
