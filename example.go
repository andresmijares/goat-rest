package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andresmijares/goat-rest/goat"
)

var (
	httpClient = customClient()
)

func customClient() goat.Client {
	clientBuilder := goat.NewBuilder()
	clientBuilder.SetMaxIdleConnections(20)
	clientBuilder.SetConnectionTimeout(10 * time.Second)
	clientBuilder.SetResponseTimeout(10 * time.Second)


	headers := make(http.Header)
	headers.Set("Authorization", "Bearer ABC1234")
	clientBuilder.SetHeaders(headers)

	client := clientBuilder.Create()

	return client
}

func main() {
	for i:= 0; i < 1000; i++ {
		go func() {
			callGithub()
		}()
	}
	time.Sleep(20*time.Second)
}

func callGithub() {
	res, err := httpClient.Get("http://github.com", nil)
	if err != nil {
		// fmt.Print(err.Error())
	} else {
		fmt.Print(res.String())
	}
}

/*
response.Body.Close()

bytes, err := ioutil.ReadAll(responbse.Body)
if err != nil {
	panic(nil)
}

var user User
if err := json.Unmarshal(bytes, &user); err != nil {
	panic(err)
}

 // user ?
*/