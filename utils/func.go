package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/gofrs/uuid"
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

func IsImage(handler *multipart.FileHeader) bool {
	typeimg := handler.Header.Get("Content-Type")
	if len(typeimg) > 6 && typeimg[:6] != "image/" {
		return false
	} else {
		return true
	}
}

func AllowRequest(token string, maxReques int, RateInterval time.Duration) bool {
	Mut.Lock()
	defer Mut.Unlock()

	now := time.Now()

	// if user not have any token in map
	info, exists := RateLimitData[token]
	if !exists {
		RateLimitData[token] = &RateLimitInfo{
			RequestCount: 1,
			LastRequest:  now,
		}
		return true
	}

	if now.Sub(info.LastRequest) > RateInterval {
		info.RequestCount = 1
		info.LastRequest = now
		return true
	}
	if info.RequestCount < maxReques {
		info.RequestCount++
		info.LastRequest = now
		return true
	}

	return false
}

func MiddePOST(fn http.HandlerFunc, show bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusNonAuthoritativeInfo)})
			return
		}
		if !AllowRequest(cookie.Value, 2, time.Second*4) {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode("you have 2 Post in 4 second")
			return
		}
		fn(w, r)
	}
}

func MiddeLike(fn http.HandlerFunc, show bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusNonAuthoritativeInfo)})
			return
		}
		if !AllowRequest(cookie.Value, 10, time.Second*2) {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(http.StatusText(http.StatusTooManyRequests))
			return
		}
		fn(w, r)
	}
}

func MiddeSingIn(fn http.HandlerFunc, show bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := Ip(r.RemoteAddr)
		if !AllowRequest(ip, 8, time.Hour*24) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("you had try all request in day , try anther day")
			return
		}
		fn(w, r)
	}
}

func MiddeSingUp(fn http.HandlerFunc, show bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := Ip(r.RemoteAddr)
		if !AllowRequest(ip, 8, time.Hour) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("you had try all request in hour , try anther later")
			return
		}
		fn(w, r)
	}
}

func Ip(str string) string {
	if strings.Contains(str, "[::1]") {
		return "[::1]"
	} else {
		ip := strings.Split(str, ":")
		return ip[0]
	}
}

func SetEnv() {
	file, err := os.Open(".env")
	if err != nil {
		log.Fatalf("Error opening .env file: %v", err)
		return
	}
	defer file.Close() // Ensure the file is closed after the function exits

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Skipping malformed line: %s", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" || value == "" {
			log.Printf("Skipping line with empty key: %s", line)
			continue
		}
		err = os.Setenv(key, value)
		if err != nil {
			log.Fatalf("Error setting environment variable: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}
}

func GenratePass() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	passwd, err := bcrypt.GenerateFromPassword([]byte(uid.String()), bcrypt.DefaultCost)
	return string(passwd), err
}
