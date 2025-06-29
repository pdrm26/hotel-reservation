package core

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(err)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func UnAuthorizedError() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized request",
	}
}

func ForbiddenError() Error {
	return Error{
		Code: http.StatusForbidden,
		Err:  "forbidden",
	}
}

func InvalidIDError() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid id given",
	}
}

func TokenExpiredError() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "token expired",
	}
}

func TokenInvalidError() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "token is invalid",
	}
}
func TokenMissingError() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "token is missing",
	}
}

func NotFoundError(resource string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  resource + "not found",
	}
}

func BadRequestError() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "bad request",
	}
}
