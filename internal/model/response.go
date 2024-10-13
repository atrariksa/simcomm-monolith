package model

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginData struct {
	Token string `json:"token"`
}
