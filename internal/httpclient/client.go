package httpclient

import (
	"net/http"
	"time"
)

// DefaultHTTPClient is the default HTTP client used by discordgo.
var DefaultHTTPClient = &http.Client{
	Jar: Cookie(),
	Transport: &http.Transport{
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 20,
	},
	Timeout: 30 * time.Second,
}
