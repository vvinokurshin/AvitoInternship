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
	"time"
)

func createConfig() *config.Config {
	cfg := new(config.Config)
	cfg.DB.DBSchemaName = "app"
	cfg.DB.DBHistoryTableName = "history"

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

func TestRepository_SelectRecordsByDate(t *testing.T) {
	cfg := createConfig()

	year := 2023
	month := time.Month(1)
	datetime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	var fakeRecords []models.History
	generateFakeData(&fakeRecords)
	fakeRecords = fakeRecords[:1]

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id", "segment_slug", "operation", "datetime"}).
		AddRow(fakeRecords[0].UserID, fakeRecords[0].SegmentSlug, fakeRecords[0].Operation, fakeRecords[0].Datetime)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "history"."user_id","history"."segment_slug","history"."operation","history"."datetime"
FROM "app"."history" WHERE date_trunc('month', datetime) = $1`)).WithArgs(datetime).WillReturnRows(rows)

	historyRep := New(cfg, gormDB)
	response, err := historyRep.SelectRecordsByDate(year, month)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeRecords, response)
	}
}
