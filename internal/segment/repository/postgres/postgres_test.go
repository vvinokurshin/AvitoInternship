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
	cfg.DB.DBU2STableName = "users2segments"
	cfg.DB.DBSegmentTableName = "segments"

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

func TestRepository_InsertSegment(t *testing.T) {
	cfg := createConfig()

	var fakeSegment *models.Segment
	generateFakeData(&fakeSegment)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	createUserRow := sqlmock.NewRows([]string{"segment_id"}).
		AddRow(fakeSegment.SegmentID)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "app"."segments" ("slug","percent","segment_id")
	VALUES ($1,$2,$3) RETURNING "segment_id"`)).WithArgs(fakeSegment.Slug, fakeSegment.Percent, fakeSegment.SegmentID).
		WillReturnRows(createUserRow)
	mock.ExpectCommit()

	segmentRep, err := New(cfg, gormDB)
	segmentID, err := segmentRep.InsertSegment(fakeSegment)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeSegment.SegmentID, segmentID)
	}
}

func TestRepository_DeleteSegment(t *testing.T) {
	cfg := createConfig()

	segmentSlug := "test"

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "app"."segments" WHERE slug = $1`)).
		WithArgs(segmentSlug).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	segmentRep, err := New(cfg, gormDB)
	err = segmentRep.DeleteSegment(segmentSlug)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_SelectUserByID(t *testing.T) {
	cfg := createConfig()

	var fakeSegment *models.Segment
	generateFakeData(&fakeSegment)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"segment_id", "slug", "percent"}).
		AddRow(fakeSegment.SegmentID, fakeSegment.Slug, fakeSegment.Percent)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "app"."segments" WHERE slug = $1`)).WithArgs(fakeSegment.Slug).WillReturnRows(rows)

	segmentRep, err := New(cfg, gormDB)
	response, err := segmentRep.SelectSegmentBySlug(fakeSegment.Slug)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeSegment, response)
	}
}

func TestRepository_SelectSegmentsByUser(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	fakeSegment := []models.Segment{
		{
			SegmentID: 1,
			Slug:      "test",
		},
	}

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"segment_id", "slug", "percent"}).
		AddRow(fakeSegment[0].SegmentID, fakeSegment[0].Slug, fakeSegment[0].Percent)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT app.segments.* FROM "app"."segments" JOIN app.users2segments using(segment_id) WHERE user_id = $1`)).
		WithArgs(userID).WillReturnRows(rows)

	segmentRep, err := New(cfg, gormDB)
	response, err := segmentRep.SelectSegmentsByUser(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeSegment, response)
	}
}

func TestRepository_InsertSegmentsToUser(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	segments := []models.AddUserToSegment{
		{
			SegmentSlug: "test",
			SegmentID:   1,
		},
	}

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "app"."users2segments" ("user_id","segment_id","until")
	VALUES ($1,$2,$3) ON CONFLICT ("user_id","segment_id") DO UPDATE SET "until"="excluded"."until"`)).
		WithArgs(userID, segments[0].SegmentID, segments[0].Until).WillReturnResult(sqlmock.NewResult(int64(0), 1))
	mock.ExpectCommit()

	segmentRep, err := New(cfg, gormDB)
	err = segmentRep.InsertSegmentsToUser(userID, segments)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_DeleteSegmentsFromUser(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	segmentIDs := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "app"."users2segments" WHERE user_id = $1 AND segment_id IN ($2)`)).
		WithArgs(userID, segmentIDs).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	segmentRep, err := New(cfg, gormDB)
	err = segmentRep.DeleteSegmentsFromUser(userID, []uint64{segmentIDs})
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestRepository_InsertUsersToSegment(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	segmentID := uint64(1)

	db, gormDB, mock, err := mockDB()
	if err != nil {
		t.Fatalf("error while mocking database: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "app"."users2segments" ("user_id","segment_id","until")
	VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`)).WithArgs(userID, segmentID, nil).
		WillReturnResult(sqlmock.NewResult(int64(0), 1))
	mock.ExpectCommit()

	segmentRep, err := New(cfg, gormDB)
	err = segmentRep.InsertUsersToSegment(segmentID, []uint64{userID})
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}
