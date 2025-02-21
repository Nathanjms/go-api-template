package application

type ResponseData map[string]interface{}

type Response struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    ResponseData        `json:"data"`
	Errors  map[string][]string `json:"errors"`
}
