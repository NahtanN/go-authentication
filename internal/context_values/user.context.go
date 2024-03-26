package context_values

type userIdContextKey string

const UserIdKey = userIdContextKey("middleware.jwt_validation.userId")
