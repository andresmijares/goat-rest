package examples

import "fmt"

type Endpoints struct {
	CurrentUserUrl string `json:"current_user_url"`
	RespositoryUrl string `json:"repository_url"`
}

func Get() (*Endpoints, error) {
	response, err := httpClient.Get("https://api.github.com", nil)
	if err != nil {
		// Manage the error as needed
		return nil, nil
	}

	fmt.Printf(fmt.Sprintf("Status Code: %d \n", response.StatusCode))
	fmt.Printf(fmt.Sprintf("Status: %s \n", response.Status))
	fmt.Printf(fmt.Sprintf("Body: %s \n", response.String()))

	var endpoints Endpoints
	if err := response.UnmarshalJson(&endpoints); err != nil {
		// Manage the error as needed
		return nil, nil
	}

	fmt.Printf(fmt.Sprintf("Current User URL: %s \n", endpoints.CurrentUserUrl))
	fmt.Printf(fmt.Sprintf("Repository URL: %s \n", endpoints.RespositoryUrl))

	return &endpoints, nil
}