package httpclient

import (
	"log"
	"net/http/cookiejar"
)

func Cookie() *cookiejar.Jar {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	return jar
}
