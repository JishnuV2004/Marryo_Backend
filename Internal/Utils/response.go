package utils

/****Add the Response JSON model structs here***/
type Success struct {
	Code   int         `json:"code"`
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

type SuccessMsg struct {
	Code    int         `json:"code"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Error struct {
	Code   int    `json:"code"`
	Status bool   `json:"status"`
	Error  string `json:"error"`
}



func SuccessResponse(data interface{}) Success {
	result := Success{
		Code:   200,
		Status: true,
		Data:   data,
	}
	return result
}

func SuccessResponseMsg(data interface{}, message string) SuccessMsg {
	result := SuccessMsg{
		Code:    200,
		Status:  true,
		Message: message,
		Data:    data,
	}
	return result
}

func ErrorMessage(code int, e error) Error {

	errMsg := Error{
		Status: false,
		Code:   code,
		Error:  e.Error(),
	}

	return errMsg
}