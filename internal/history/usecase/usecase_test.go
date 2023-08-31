package usecase

import (
	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	mockHistoryRepo "github.com/vvinokurshin/AvitoInternship/internal/history/repository/mocks"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	"github.com/vvinokurshin/AvitoInternship/pkg"
	"testing"
)

func createConfig() *config.Config {
	return new(config.Config)
}

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestUseCase_GetHistoryCSV(t *testing.T) {
	pkg.HistoryFolderName = "../../../history/"
	cfg := createConfig()

	fakeForm := models.FormHistory{
		Year:  2023,
		Month: 1,
	}
	var fakeHistoryResponse []models.History
	generateFakeData(&fakeHistoryResponse)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	historyRepo := mockHistoryRepo.NewMockRepositoryI(ctrl)
	historyUC := New(cfg, historyRepo)

	historyRepo.EXPECT().GetRecordsByDate(fakeForm.Year, fakeForm.Month).Return(fakeHistoryResponse, nil)
	response, err := historyUC.GetHistoryCSV(fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, "history-2023-1.csv", response)
	}
}
