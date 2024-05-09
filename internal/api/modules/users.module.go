package modules

import (
	"github.com/nahtann/go-lab/internal/api/router"
	"github.com/nahtann/go-lab/internal/handlers/users_handlers/current_user"
	"github.com/nahtann/go-lab/internal/middlewares"
	"github.com/nahtann/go-lab/internal/storage/database/models"
	"github.com/nahtann/go-lab/internal/utils"
	"github.com/nahtann/go-lab/internal/wrappers"
)

const usersRootRoute = "/users"

func UsersModule(router *router.ApiRouter) {
	currentUserRoute(router)
}

// @Summary		Get current user data.
// @Description	Must be authenticated
// @Tags			users
// @Produce		json
// @Success		200	{object}	models.UserModel
// @router			/users/current [get]
// @Security		ApiKeyAuth
func currentUserRoute(router *router.ApiRouter) {
	handler := current_user.Handler{
		DB: router.DB,
	}

	httpWrapper := wrappers.HttpWrapper[current_user.Request, models.UserModel]{
		Handler: &handler,
		RequestParsers: []wrappers.RequestParser[current_user.Request]{
			current_user.RequestParser,
		},
		ValidateRequest: utils.Validate,
	}

	router.SetRoute(
		"GET",
		utils.SetSubRoute(usersRootRoute, "/current"),
		&httpWrapper,
		middlewares.NewJWTValidationHttpHandler(),
	)
}
