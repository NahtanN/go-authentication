package current_user

/*import (*/
/*"fmt"*/
/*"testing"*/
/*"time"*/

/*"github.com/pashagolub/pgxmock/v3"*/

/*"github.com/nahtann/go-lab/internal/storage/database/models"*/
/*"github.com/nahtann/go-lab/internal/utils"*/
/*"github.com/nahtann/go-lab/internal/utils/mocks"*/
/*)*/

/*func TestShouldGetCurrentUser(t *testing.T) {*/
/*user := models.UserModel{*/
/*Id:        1,*/
/*Username:  "Test User",*/
/*Email:     "test@test.com",*/
/*CreatedAt: time.Now(),*/
/*}*/

/*mock, err := pgxmock.NewPool(*/
/*pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual),*/
/*)*/
/*if err != nil {*/
/*t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)*/
/*}*/
/*defer mock.Close()*/

/*rows := mock.NewRows([]string{"username", "email", "created_at"}).*/
/*AddRow(user.Username, user.Email, user.CreatedAt)*/

/*mock.ExpectQuery("SELECT username, email, created_at FROM users WHERE id = $1").*/
/*WithArgs(user.Id).*/
/*WillReturnRows(rows)*/

/*responseUser, err := CurrentUser(mock, user.Id)*/
/*if err != nil {*/
/*t.Fatalf("an error '%s' was not expected when trying to get current user", err)*/
/*}*/

/*mocks.AssertJSON(responseUser, user, t)*/
/*}*/

/*func TestShouldFailOnParseDbData(t *testing.T) {*/
/*user := models.UserModel{*/
/*Id:       1,*/
/*Username: "Test User",*/
/*}*/

/*mock, err := pgxmock.NewPool(*/
/*pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual),*/
/*)*/
/*if err != nil {*/
/*t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)*/
/*}*/
/*defer mock.Close()*/

/*rows := mock.NewRows([]string{"1"}).*/
/*AddRow(user.Username)*/

/*mock.ExpectQuery("SELECT username, email, created_at FROM users WHERE id = $1").*/
/*WithArgs(user.Id).*/
/*WillReturnRows(rows)*/

/*_, err = CurrentUser(mock, user.Id)*/
/*if err == nil {*/
/*t.Fatalf("Does not triggered error")*/
/*}*/

/*expected := utils.CustomError{*/
/*Message: "Unable to parse current user data.",*/
/*}*/

/*mocks.AssertJSON(err, expected, t)*/
/*}*/

/*func TestShouldFailOnDbError(t *testing.T) {*/
/*mock, err := pgxmock.NewPool(*/
/*pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual),*/
/*)*/
/*if err != nil {*/
/*t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)*/
/*}*/
/*defer mock.Close()*/

/*mock.ExpectQuery("SELECT username, email, created_at FROM users WHERE id = $1").*/
/*WillReturnError(fmt.Errorf("Some error"))*/

/*_, err = CurrentUser(mock, 1)*/

/*if err == nil {*/
/*t.Fatalf("Does not triggered error")*/
/*}*/

/*expected := utils.CustomError{*/
/*Message: "Unable to retrieve current user data.",*/
/*}*/

/*mocks.AssertJSON(err, expected, t)*/
/*}*/
