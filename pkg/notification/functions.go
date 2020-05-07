package notification

import (
	"log"
	"net/http"
	"strings"
)

type Func func(msg string)

func HttpFunc(url string) Func {
	return func(message string) {
		response, err := http.Post(url, "text/plain", strings.NewReader(message))
		if err != nil {
			log.Printf("сan't send notification, error is %s", err)
		} else if (response.StatusCode/100)*100 != http.StatusOK {
			log.Printf("got status code %d while sending notification", response.StatusCode)
		}
	}
}
