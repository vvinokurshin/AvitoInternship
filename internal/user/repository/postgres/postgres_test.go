package postgres

import (
	"database/sql"
	"fmt"
	"github.com/go-faker/faker/v4"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func createConfig() *config.Config {
	cfg := new(config.Config)
	cfg.DB.DBSchemaName = "app"
	cfg.DB.DBUserTableName = "users"

	return cfg
}

func mockDB() (*sql.DB, *gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("mocking database error: %s", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("opening gorm error: %s", err)
	}

	return db, gormDB, mock, nil
}

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestRepository_InsertUser(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	createUserRow := sqlmock.NewRows([]string{"user_id"}).
		AddRow(fakeUser.UserID)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "app"."users" ("username","first_name","last_name","user_id")
	VALUES ($1,$2,$3,$4) RETURNING "user_id"`)).WithArgs(fakeUser.Username, fakeUser.FirstName, fakeUser.LastName, fakeUser.UserID).
		WillReturnRows(createUserRow)
	mock.ExpectCommit()

	userRep := New(cfg, gormDB)
	userID, err := userRep.InsertUser(fakeUser)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUser.UserID, userID)
	}
}

func TestRepository_UpdateUser(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "app"."users" SET "username"=$1,"first_name"=$2,"last_name"=$3 WHERE "user_id" = $4`)).
		WithArgs(fakeUser.Username, fakeUser.FirstName, fakeUser.LastName, fakeUser.UserID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	userRep := New(cfg, gormDB)
	err = userRep.UpdateUser(fakeUser)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_DeleteUser(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "app"."users" WHERE user_id = $1`)).
		WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	userRep := New(cfg, gormDB)
	err = userRep.DeleteUser(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_SelectUserByID(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id", "username", "first_name", "last_name"}).
		AddRow(fakeUser.UserID, fakeUser.Username, fakeUser.FirstName, fakeUser.LastName)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "app"."users" WHERE user_id = $1`)).WithArgs(fakeUser.UserID).WillReturnRows(rows)

	userRep := New(cfg, gormDB)
	response, err := userRep.SelectUserByID(fakeUser.UserID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}

func TestRepository_SelectUserByUsername(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id", "username", "first_name", "last_name"}).
		AddRow(fakeUser.UserID, fakeUser.Username, fakeUser.FirstName, fakeUser.LastName)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "app"."users" WHERE username = $1`)).WithArgs(fakeUser.Username).WillReturnRows(rows)

	userRep := New(cfg, gormDB)
	response, err := userRep.SelectUserByUsername(fakeUser.Username)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}

func TestRepository_GetUserIDs(t *testing.T) {
	cfg := createConfig()

	fakeUserIDs := []uint64{1}

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id"}).AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT user_id FROM "app"."users"`)).WillReturnRows(rows)

	userRep := New(cfg, gormDB)
	response, err := userRep.SelectUserIDs()
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUserIDs, response)
	}
}
