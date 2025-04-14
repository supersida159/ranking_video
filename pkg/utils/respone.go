package helper

import (
	apperror "ranking_video/pkg/app_error"
	"time"

	"github.com/gin-gonic/gin"
)

type EmptyObj struct{}
type Response struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Errors    interface{} `json:"errors,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	TimeStamp time.Time   `json:"time_stamp"`
}

func BuildSuccessResponse(status int, message string, data interface{}) Response {
	res := Response{
		Status:    status,
		Message:   message,
		Errors:    nil,
		Data:      data,
		TimeStamp: time.Now(),
	}
	return res
}

func BuildErrorResponse(appErr *apperror.AppError) Response {
	res := Response{
		Status:    int(appErr.StatusCode),
		Message:   "",
		Errors:    appErr,
		Data:      nil,
		TimeStamp: time.Now(),
	}
	return res
}

func BuildSuccessGinResponse(c *gin.Context, data interface{}, message ...string) {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = "Success" // default message if none provided
	}

	res := Response{
		Status:    c.Writer.Status(), // get status from the context
		Message:   msg,
		Errors:    nil,
		Data:      data,
		TimeStamp: time.Now(),
	}

	c.JSON(c.Writer.Status(), res) // status comes from the context writer
}

func BuildErrorGinResponse(c *gin.Context, appErr *apperror.AppError) {
	res := Response{
		Status:    int(appErr.StatusCode),
		Message:   appErr.MessageEn, // Assuming you'd want to send the English message, modify as needed
		Errors:    appErr,
		Data:      nil,
		TimeStamp: time.Now(),
	}
	c.JSON(appErr.StatusCode, res)

}

func BuildErrorGinResponseAndAbort(c *gin.Context, appErr *apperror.AppError) {
	res := Response{
		Status:    int(appErr.StatusCode),
		Message:   appErr.MessageEn, // Assuming you'd want to send the English message, modify as needed
		Errors:    appErr,
		Data:      nil,
		TimeStamp: time.Now(),
	}
	c.JSON(appErr.StatusCode, res)
	c.Abort()
}
