package examples

import (
	"net/http"
	"time"

	"github.com/andresmijares/goat-rest/goat"
)

var (
	httpClient = getHtppClient()
)

func getHtppClient() goat.Client {
	client := goat.New().
		SetConnectionTimeout(2 * time.Second).
		SetResponseTimeout(3 * time.Second). 
		Create()
	return client
}

func CustomHttpClient() goat.Client {
	currentClient := http.Client{}
	client := goat.New().
		SetConnectionTimeout(2 * time.Second).
		SetResponseTimeout(3 * time.Second). 
		SetHttpClient(&currentClient).
		Create()

	return client
}