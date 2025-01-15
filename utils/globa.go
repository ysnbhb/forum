package utils

import "sync"

var (
	GithubReURL   = "https://localhost:8080/auth/github/callback"
	GoogleReURL   = "https://localhost:8080/auth/google/callback"
	RateLimitData = make(map[string]*RateLimitInfo)
	Mut           sync.Mutex
)
