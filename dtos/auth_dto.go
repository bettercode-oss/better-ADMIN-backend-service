package dtos

import "github.com/labstack/echo"

type MemberSignIn struct {
	Id       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (m MemberSignIn) Validate(ctx echo.Context) error {
	return ctx.Validate(m)
}
