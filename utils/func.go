package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
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

	tmpl, err := template.ParseFiles("./view/error.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, msg1, statusCode)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, Error); err != nil {
		fmt.Println(err)
		http.Error(w, msg1, statusCode)
		return
	}
	w.WriteHeader(statusCode)
	// If successful, write the buffer content to the ResponseWriter
	buf.WriteTo(w)
}

func IsValidUsername(username string) bool {
	if username == "" || len(username) > 10 {
		return false
	}
	last := []rune(username)[0]
	for _, c := range username {
		if c == '_' {
			if last == '_' {
				return false
			}
			last = c
			continue
		}
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
		last = c
	}
	return true
}

// This for valid email

func IsValidEmail(email string) bool {
	re := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(re)
	return regex.MatchString(email)
}

func HasPassowd(password string) (string, error) {
	hashpassord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashpassord), nil
}

func IsValidPassword(password string) bool {
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, ch := range password {
		if unicode.IsUpper(ch) {
			hasUpper = true
		} else if unicode.IsLower(ch) {
			hasLower = true
		} else if unicode.IsDigit(ch) {
			hasDigit = true
		} else if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func IsImage(filePath io.Reader) (bool, error) {
	buffer := make([]byte, 512)
	_, err := filePath.Read(buffer)
	if err != nil {
		return false, err
	}
	contentType := http.DetectContentType(buffer)
	fmt.Println("Content Type:", contentType)
	return contentType == "image/jpeg" ||
		contentType == "image/png" ||
		contentType == "image/gif" ||
		contentType == "image/webp" ||
		contentType == "image/svg+xml", nil
}
