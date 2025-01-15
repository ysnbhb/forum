package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"forum/database"
	"forum/server"
	"forum/utils"
)

func main() {
	DB := database.IntDB()
	err := database.CreateTable(DB)
	if err != nil {
		log.Fatal("Error creating tables: ", err)
		return
	}
	mux := http.NewServeMux()
	api := server.New(DB)
	mux.HandleFunc("/signup", (api.PageSingUp))
	mux.HandleFunc("/", api.Page)
	mux.HandleFunc("/signin", (api.PageSingIn))
	mux.HandleFunc("/logout", utils.MiddePOST(api.LogOut, false))
	mux.HandleFunc("/user/signup", utils.MiddeSingUp(api.DB.SingUp, false))
	mux.HandleFunc("/user/signin", utils.MiddeSingIn(api.DB.SingIn, false))
	mux.HandleFunc("/user/exist", (api.DB.Exist))
	mux.HandleFunc("/post/lastId", api.DB.LastID)
	mux.HandleFunc("/frontend/", server.Server)
	mux.HandleFunc("/userImg/", server.Server)
	mux.HandleFunc("/auth/google", DB.GoogleAthud)
	mux.HandleFunc("/auth/google/callback", DB.GoogleCallbackHandler)
	mux.HandleFunc("/auth/github", DB.GithubLoginHandler)
	mux.HandleFunc("/auth/github/callback", DB.GithubCallbackHandler)
	// api
	mux.HandleFunc("/api/getCategorie", api.DB.GetCtg)
	mux.HandleFunc("/api/posts", (DB.AllPost))
	mux.HandleFunc("/api/post", (api.DB.OnePost))
	mux.HandleFunc("/api/addPost", utils.MiddePOST(api.DB.AddPost, false))
	mux.HandleFunc("/api/post/like", utils.MiddeLike(DB.LikePost, false))
	mux.HandleFunc("/api/commant/like", utils.MiddeLike(DB.LikeCommat, false))
	mux.HandleFunc("/api/commant", (DB.GetCommat))
	mux.HandleFunc("/api/addcommat", utils.MiddePOST(DB.AddCommat, false))
	mux.HandleFunc("/api/commatedPost", DB.FilterCommatedPost)
	mux.HandleFunc("/api/notif", DB.Notif)
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
	utils.SetEnv()
	port := os.Getenv("PORT")
	cert := os.Getenv("PuplcKey")
	key := os.Getenv("PriveKey")
	if key == "" || port == "" || cert == "" {
		log.Fatal("Set Env For PORT || PuplcKey || PriveKey")
	}
	srv := &http.Server{
		Addr:         port,
		Handler:      mux,
		TLSConfig:    tlsConfig,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}
	// for generate cert file use this commad go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost in golang
	fmt.Printf("server is run in https://localhost%s/\n", port)
	err = srv.ListenAndServeTLS(cert, key)
	// err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
