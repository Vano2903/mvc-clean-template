package httpserver

import "github.com/labstack/echo/v4"

// ConnectionID string `json:"connection_id,omitemtpy"`
type HttpError struct {
	Code        int    `json:"code" example:"400"`
	IsError     bool   `json:"is_error" example:"true"`
	Message     string `json:"message" example:"Bad Request"`
	Details     string `json:"details,omitempty" example:"Bad Request With More Info"`
	Instance    string `json:"instance,omitempty" example:"/api/v1/users/1"`
	ErrorType   string `json:"error_type,omitempty" example:"invalid_id"`
	ErrorDocUrl string `json:"error_doc_url,omitempty" example:"https://example.com/docs/errors/invalid_id"`
}

func respError(c echo.Context, code int, message, details, errType string) error {
	h := HttpError{
		Instance:    c.Request().RequestURI,
		IsError:     true,
		Code:        code,
		Message:     message,
		Details:     details,
		ErrorType:   errType,
		ErrorDocUrl: "https://example.com/docs/errors/" + errType, //could be a map somwhere, just an example for now
	}

	return c.JSON(code, h)
}

type HttpSuccess struct {
	Code    int         `json:"code" example:"200"`
	IsError bool        `json:"is_error" example:"false"`
	Message string      `json:"message" example:"OK"`
	Data    interface{} `json:"data,omitempty"`
}

func respSuccess(c echo.Context, code int, message string, data ...interface{}) error {
	h := HttpSuccess{
		Code:    code,
		IsError: false,
		Message: message,
	}

	if len(data) > 0 {
		h.Data = data[0]
	}

	return c.JSON(code, h)
}
