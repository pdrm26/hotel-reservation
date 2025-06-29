package core

import "github.com/gofiber/fiber/v2"

func ErrorHandler(c *fiber.Ctx, err error) error {
	return c.JSON(map[string]string{"error": err.Error()})
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

func NotAuthorizedError() Error {
	return Error{Code: 401, Err: "unauthorized"}
}
