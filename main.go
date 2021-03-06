package main

import(
	"net/http"
	"os"
	"io"
	"fmt"
)

func uploadsHandler(w http.ResponseWriter, r *http.Request){
	uploadFile, header, err := r.FormFile("upload_file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w,err)
		return
	}
	defer uploadFile.Close()
	dirname := "./uploads"
	os.MkdirAll(dirname,0755)
	filepath := fmt.Sprintf("%s/%s",dirname,header.Filename)
	file, err1 := os.Create(filepath)
	defer file.Close()
	if err1 != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w,err)
		return
	}
	io.Copy(file,uploadFile)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w,filepath)
}

func main(){
	http.HandleFunc("/uploads",uploadsHandler)
	http.Handle("/",http.FileServer(http.Dir("public")))
	
	http.ListenAndServe(":3000",nil)
}
