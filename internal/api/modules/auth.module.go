package modules

import (
	"github.com/nahtann/go-lab/internal/api/router"
	"github.com/nahtann/go-lab/internal/handlers/auth_handlers/sign_in"
	"github.com/nahtann/go-lab/internal/interfaces"
	"github.com/nahtann/go-lab/internal/utils"
	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

const authRootRoute = "/auth"

func AuthModule(router *router.ApiRouter) {
	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/sign-in"),
		signRoute(router.DB),
	)
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

func signRoute(db interfaces.Pgx) *sign_in.HttpWrapper {
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
