package api

// BaseResponse denotes the base response structure for the API response
type BaseResponse struct {
	Data          interface{} `json:"data"`
	StatusMessage string      `json:"status_message"`
	StatusCode    int         `json:"status_code"`
}
