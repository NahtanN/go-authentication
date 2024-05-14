package modules

import (
	"github.com/nahtann/go-lab/internal/api/router"
	"github.com/nahtann/go-lab/internal/handlers/auth_handlers/refresh_token"
	"github.com/nahtann/go-lab/internal/handlers/auth_handlers/sign_in"
	"github.com/nahtann/go-lab/internal/handlers/auth_handlers/sign_up"
	"github.com/nahtann/go-lab/internal/middlewares"
	"github.com/nahtann/go-lab/internal/utils"
	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
	wrapper_utils "github.com/nahtann/go-lab/internal/utils/wrapper"
	"github.com/nahtann/go-lab/internal/wrappers"
)

const authRootRoute = "/auth"

func AuthModule(router *router.ApiRouter) {
	signInRoute(router)
	signUpRoute(router)
	refreshTokenRoute(router)
}

// @Description	Authenticate user and returns access and refresh tokens.
// @Tags			auth
// @Accept			json
// @Param			request	body	SigninRequest	true	"Request Body"
// @Produce		json
// @Success		201	{object}	auth_utils.Tokens
// @Failure		400	{object}	utils.CustomError	"Message: 'User or password invalid.'"
// @router			/auth/sign-in   [post]
func signInRoute(router *router.ApiRouter) {
	signIn := sign_in.Handler{
		DB:              router.DB,
		VerifyPassword:  utils.VerifyPassword,
		CreateJwtTokens: auth_utils.CreateJwtTokens,
	}

	httpWrapper := wrappers.HttpWrapper[sign_in.Request, auth_utils.Tokens]{
		Handler: &signIn,
		RequestParsers: []wrappers.RequestParser[sign_in.Request]{
			wrapper_utils.BodyParser[sign_in.Request],
		},
		ValidateRequest: utils.Validate,
	}

	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/sign-in"),
		&httpWrapper,
	)
}

// @Description	Creates new user.
// @Tags			auth
// @Accept			json
// @Param			request	body	sign_up.SignupRequest	true	"Request Body"
// @Produce		json
// @Success		201	{object}	utils.DefaultResponse	"Message: 'Sign up successfully'"
// @Failure		400	{object}	utils.CustomError		"Message: 'Username already in use. E-mail already in use.'"
// @router			/auth/sign-up   [post]
func signUpRoute(router *router.ApiRouter) {
	signUp := sign_up.Handler{
		DB:           router.DB,
		HashPassword: utils.HashPassword,
	}

	httpWrapper := wrappers.HttpWrapper[sign_up.Request, utils.DefaultResponse]{
		Handler: &signUp,
		RequestParsers: []wrappers.RequestParser[sign_up.Request]{
			wrapper_utils.BodyParser[sign_up.Request],
		},
		ValidateRequest: utils.Validate,
	}

	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/sign-up"),
		&httpWrapper,
	)
}

// @Description	Creates a new pair of access and refresh tokens.
// @Tags			auth
// @Accept			json
// @Param			request	body	RefreshTokenRequest	true	"Request Body"
// @Produce		json
// @Success		201	{object}	auth_utils.Tokens
// @Failure		401	{object}	utils.CustomError	"Message: 'Invalid Request'"
// @router			/auth/refresh-token [post]
func refreshTokenRoute(router *router.ApiRouter) {
	invalidate := refresh_token.InvalidateHandler{
		DB: router.DB,
	}

	update := refresh_token.UpdateHandler{
		DB:           router.DB,
		CreateTokens: auth_utils.CreateJwtTokens,
	}

	refreshToken := refresh_token.Handler{
		DB:                     router.DB,
		ValidateToken:          middlewares.ValidateJWT,
		InvalidateTokensByUser: invalidate.TokensByUser,
		UpdateUserTokens:       update.UserTokens,
	}

	httpWrapper := wrappers.HttpWrapper[refresh_token.Request, auth_utils.Tokens]{
		Handler: &refreshToken,
		RequestParsers: []wrappers.RequestParser[refresh_token.Request]{
			wrapper_utils.BodyParser[refresh_token.Request],
		},
		ValidateRequest: utils.Validate,
	}

	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/refresh-token"),
		&httpWrapper,
	)
}
