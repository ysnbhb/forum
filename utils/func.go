package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, statusCode int, msg1, msg2 string, err error) {
	// print errors in case of intenal server error
	if err != nil && statusCode == 500 {
		log.Println(err)
	}

	Error := ErrorData{
		Msg1:       msg1,
		Msg2:       msg2,
		StatusCode: statusCode,
	}

	tmpl, err := template.ParseFiles("./veiw/error.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, msg1, statusCode)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, Error); err != nil {
		http.Error(w, msg1, statusCode)
		return
	}
	w.WriteHeader(statusCode)
	// If successful, write the buffer content to the ResponseWriter
	buf.WriteTo(w)
}
