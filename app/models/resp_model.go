package models

type ApiData struct {
	IsSuccess  bool                `json:"isSuccess"`
	StatusCode int                 `json:"statusCode"`
	Error      interface{}         `json:"error"`
	Result     interface{}         `json:"result"`
	Headers    []map[string]string `json:"headers"`
}

type DatabaseResult struct {
	Error  error
	Result interface{}
}
