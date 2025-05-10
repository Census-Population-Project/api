package response

type SuccessResponse struct {
	Status string `json:"status" default:"success"`
}

type SuccessResponseWithResult struct {
	Status string      `json:"status" default:"success"`
	Result interface{} `json:"result"`
}

type ErrorResponse struct {
	Status  string `json:"status" default:"error"`
	Message string `json:"message"`
}
