package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"forum/utils"

	"github.com/gofrs/uuid/v5"
)

func (db *Date) GoogleAthud(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err == nil {
		if db.TakeIDfromAuth(cookie.Value) {
			http.Redirect(w, r, "/auth/signup", http.StatusFound)
		}
	}
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	authURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid%%20email",
		googleClientID,
		utils.GoogleReURL,
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (db *Date) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/sign", http.StatusSeeOther)
		return
	}
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	data := url.Values{}
	data.Set("client_id", googleClientID)
	data.Set("client_secret", googleClientSecret)
	data.Set("redirect_uri", utils.GoogleReURL)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&tokenResp)

	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	userResp, err := http.DefaultClient.Do(req)
	if err != nil || userResp.StatusCode != http.StatusOK {
		http.Redirect(w, r, "/sign", http.StatusSeeOther)
		return
	}
	defer userResp.Body.Close()
	Google := utils.Google{}
	json.NewDecoder(userResp.Body).Decode(&Google)
	db.HandleAuth(Google.Email, w, r)
}

func (db *Date) GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err == nil {
		if db.TakeIDfromAuth(cookie.Value) {
			http.Redirect(w, r, "/auth/signup", http.StatusFound)
		}
	}
	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	authURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user",
		githubClientID,
		(utils.GithubReURL),
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// GitHub Callback
func (db *Date) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/sign", http.StatusSeeOther)
		return
	}

	// Exchange code for token
	data := url.Values{}
	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	data.Set("client_id", githubClientID)
	data.Set("client_secret", githubClientSecret)
	data.Set("code", code)

	resp, err := http.PostForm("https://github.com/login/oauth/access_token", data)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Redirect(w, r, "/sign", http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	tokenResp, _ := url.ParseQuery(string(body))
	accessToken := tokenResp.Get("access_token")

	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	userResp, err := http.DefaultClient.Do(req)
	if err != nil || userResp.StatusCode != http.StatusOK {
		http.Redirect(w, r, "/sign", http.StatusSeeOther)
		return
	}
	defer userResp.Body.Close()
	gitub := utils.GitHub{}
	json.NewDecoder(userResp.Body).Decode(&gitub)
	db.HandleAuth(gitub.Login+"@github.com", w, r)
}

func (db *Date) InsertAuth(email string) error {
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	query := `INSERT INTO noLog(email , token) VALUES (? , ?)`
	db.DB.Exec(query, email, uid.String())
	return nil
}

func (db *Date) HandleAuth(email string, w http.ResponseWriter, r *http.Request) {
	id, typeOflog := db.SelectAuth(email)
	if id == -1 {
		uid, _ := uuid.NewV4()
		query := `INSERT INTO noLog (email , token) VALUES(?, ?)`
		db.DB.Exec(query, email, uid)
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  uid.String(),
			Path:   "/",
			MaxAge: 60,
		})
		http.Redirect(w, r, "/auth/sigup", http.StatusSeeOther)
		return
	}
	if typeOflog == "sing" {
		http.Redirect(w, r, "/auth/used", http.StatusSeeOther)
	} else {
		db.SetCookie(w, id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
