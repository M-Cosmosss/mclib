package route

import (
	"github.com/M-Cosmosss/mclib/internal/context"
	"net/http"
)

func userAuth(ctx *context.Context) error {
	if !ctx.IsUser {
		return ctx.Error(http.StatusUnauthorized, "用户未登录")
	}
	return nil
}

func managerAuth(ctx *context.Context) error {
	if !ctx.IsManager {
		return ctx.Error(http.StatusUnauthorized, "权限不足")
	}
	return nil
}
