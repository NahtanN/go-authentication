package refresh_token

import (
	"fmt"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestInvalidateHandler(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer mock.Close()

	userId := uint32(1)

	mock.ExpectExec("UPDATE refresh_tokens SET used = true WHERE user_id = $1").
		WithArgs(userId).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	handler := InvalidateHandler{
		DB: mock,
	}

	err = handler.TokensByUser(userId)

	assert.Nil(t, err)
}

func TestInvalidateHandlerQueryError(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer mock.Close()

	userId := uint32(1)

	mock.ExpectExec("UPDATE refresh_tokens SET used = true WHERE user_id = $1").
		WithArgs(userId).
		WillReturnError(fmt.Errorf("Some error"))

	handler := InvalidateHandler{
		DB: mock,
	}

	err = handler.TokensByUser(userId)

	assert.NotNil(t, err)
	assert.Equal(t, "Invalid Request.", err.Error())
}
