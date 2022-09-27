package response

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ReplyOK(c *gin.Context, data interface{}) {
	HTTPComplete(c, http.StatusOK, NewResponseOk(data))
}

func ReplyError(c *gin.Context, err interface{}) {
	ReplyErrorWithData(c, err, nil)
}

func ReplyErrorWithData(c *gin.Context, err interface{}, data interface{}) {
	switch e := err.(type) {
	case *AuthorizationError:
		res := NewResponseError(http.StatusUnauthorized, e.Error(), data)
		HTTPComplete(c, http.StatusUnauthorized, res)
		break

	case validator.ValidationErrors:
		res := NewResponseError(StatusInvalidParam, e[0].Error(), data)
		HTTPComplete(c, http.StatusBadRequest, res)
		break

	case *BusinessError:
		res := NewResponseError(e.Code, e.Message, data)
		HTTPComplete(c, http.StatusOK, res)
		break

	case *IllegalArgumentError:
		res := NewResponseError(400, e.Error(), data)
		HTTPComplete(c, http.StatusBadRequest, res)
		break

	case error:
		res := NewResponseError(StatusUnknownError, e.Error(), data)
		HTTPComplete(c, http.StatusInternalServerError, res)
		break

	default:
		res := NewResponseError(StatusUnknownError, "unknown error", data)
		HTTPComplete(c, http.StatusInternalServerError, res)
	}
}

func HTTPComplete(c *gin.Context, httpStatus int, body interface{}) {
	c.JSON(httpStatus, body)
}
