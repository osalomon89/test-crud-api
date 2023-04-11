package presenter

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type UserResponse struct {
	Response
	Data *jsonUser `json:"data"`
}

type TokenResponse struct {
	Response
	AccessToken string `json:"access_token"`
}

type ItemResponse struct {
	Response
	Data *jsonItem `json:"data"`
}
