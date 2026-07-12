package middleware

import (
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"golang-clean-architecture/internal/util"

	"github.com/gofiber/fiber/v2"
)

func NewAuthAdmin(userUseCase *usecase.UserUseCase, tokenUtil *util.TokenUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		userUseCase.Log.Debugf("Authorization : %s", request.Token)

		auth, err := tokenUtil.ParseToken(ctx.UserContext(), request.Token)
		if err != nil {
			userUseCase.Log.Warnf("Failed parse token admin by token : %+v", err)
			return fiber.ErrUnauthorized
		}
		ctx.Locals("auth", auth)

		err = tokenUtil.ValidateJWT(ctx)
		if err != nil {
			userUseCase.Log.Warnf("Failed find  admin by token : %+v", err)
			return fiber.ErrUnauthorized
		}

		error := userUseCase.TokenUtil.ValidateAdminRoleJWT(ctx)
		if error != nil {
			userUseCase.Log.Warnf("Only Administrator is allowed to perform this action : %+v", err)
			return fiber.ErrUnauthorized
		}
		return ctx.Next()
	}
}

func NewAuthCustomer(userUseCase *usecase.UserUseCase, tokenUtil *util.TokenUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		userUseCase.Log.Debugf("Authorization : %s", request.Token)
		auth, err := tokenUtil.ParseToken(ctx.UserContext(), request.Token)
		if err != nil {
			userUseCase.Log.Warnf("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}
		ctx.Locals("auth", auth)
		err = tokenUtil.ValidateJWT(ctx)
		if err != nil {
			userUseCase.Log.Warnf("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}

		error := tokenUtil.ValidateCustomerRoleJWT(ctx)
		if error != nil {
			userUseCase.Log.Warnf("Only registered Customers are allowed to perform this action : %+v", err)
			return fiber.ErrUnauthorized
		}
		return ctx.Next()
	}

}

func NewAuthSuperAdmin(userUseCase *usecase.UserUseCase, tokenUtil *util.TokenUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		userUseCase.Log.Debugf("Authorization : %s", request.Token)
		auth, err := tokenUtil.ParseToken(ctx.UserContext(), request.Token)
		if err != nil {
			userUseCase.Log.Warnf("Failed find Superuser by token : %+v", err)
			return fiber.ErrUnauthorized
		}
		ctx.Locals("auth", auth)
		err = tokenUtil.ValidateJWT(ctx)
		if err != nil {
			userUseCase.Log.Warnf("Failed find Super admin by token : %+v", err)
			return fiber.ErrUnauthorized
		}
		error := tokenUtil.ValidateSuperAdminRoleJWT(ctx)
		if error != nil {
			userUseCase.Log.Warnf("Only Super Administrator is allowed to perform this action : %+v", err)
			return fiber.ErrUnauthorized
		}
		return ctx.Next()
	}
}

func NewAuthDriver(userUseCase *usecase.UserUseCase, tokenUtil *util.TokenUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		userUseCase.Log.Debugf("Authorization : %s", request.Token)
		auth, err := tokenUtil.ParseToken(ctx.UserContext(), request.Token)
		if err != nil {
			userUseCase.Log.Warnf("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}
		ctx.Locals("auth", auth)
		err = tokenUtil.ValidateJWT(ctx)
		if err != nil {
			userUseCase.Log.Warnf("Failed Auth driver user by token : %+v", err)
			return fiber.ErrUnauthorized
		}
		error := tokenUtil.ValidateDriverRoleJWT(ctx)
		if error != nil {
			userUseCase.Log.Warnf("Driver is allowed to perform this action : %+v", err)
			return fiber.ErrUnauthorized
		}
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
