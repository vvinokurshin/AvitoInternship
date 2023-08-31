package delivery

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	mockHistoryUC "github.com/vvinokurshin/AvitoInternship/internal/history/usecase/mocks"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	"github.com/vvinokurshin/AvitoInternship/pkg"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func createConfig() *config.Config {
	return new(config.Config)
}

func TestDelivery_GetHistoryCSV(t *testing.T) {
	pkg.HistoryFolderName = "../../../history/"
	cfg := createConfig()

	fakeForm := models.FormHistory{
		Year:  2023,
		Month: 1,
	}
	fakeFilename := "history-2023-1.csv"
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	historyUC := mockHistoryUC.NewMockUseCaseI(ctrl)
	historyH := New(cfg, historyUC)

	r := httptest.NewRequest(http.MethodGet, "/history", bytes.NewReader([]byte{}))
	q := r.URL.Query()
	q.Set("year", strconv.Itoa(fakeForm.Year))
	q.Set("month", strconv.Itoa(int(fakeForm.Month)))
	r.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()

	historyUC.EXPECT().GetHistoryCSV(fakeForm).Return(fakeFilename, nil)
	historyH.GetHistoryCSV(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected status %d, got %d ", status, w.Code)
	}
}
