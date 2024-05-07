package modules

import (
	"github.com/nahtann/go-lab/internal/api/router"
	"github.com/nahtann/go-lab/internal/handlers/auth_handlers/sign_in"
	"github.com/nahtann/go-lab/internal/handlers/auth_handlers/sign_up"
	"github.com/nahtann/go-lab/internal/interfaces"
	"github.com/nahtann/go-lab/internal/utils"
	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
	"github.com/nahtann/go-lab/internal/wrappers"
)

const authRootRoute = "/auth"

func AuthModule(router *router.ApiRouter) {
	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/sign-in"),
		signInRoute(router.DB),
	)

	signUpRoute(router)
	/* router.SetRoute(*/
	/*"POST",*/
	/*utils.SetSubRoute(authRootRoute, "/sign-up"),*/
	/*auth_handlers.NewSignUpHttpHandler(router.DB),*/
	/*)*/
	/*router.SetRoute(*/
	/*"POST",*/
	/*utils.SetSubRoute(authRootRoute, "/refresh-token"),*/
	/*auth_handlers.NewRefreshTokenHttpHandler(router.DB),*/
	/*)*/
}

func signInRoute(db interfaces.Pgx) *sign_in.HttpWrapper {
	signIn := sign_in.Handler{
		DB:              db,
		VerifyPassword:  utils.VerifyPassword,
		CreateJwtTokens: auth_utils.CreateJwtTokens,
	}

	signInHttpWrapper := sign_in.HttpWrapper{
		Handler:         &signIn,
		ValidateRequest: utils.Validate,
	}

	return &signInHttpWrapper
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
		Handler:         &signUp,
		ValidateRequest: utils.Validate,
	}

	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/sign-up"),
		&httpWrapper,
	)
}
