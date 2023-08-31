package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	mockSegmentUC "github.com/vvinokurshin/AvitoInternship/internal/segment/usecase/mocks"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestDelivery_CreateSegment(t *testing.T) {
	cfg := createConfig()

	var fakeForm models.FormSegment
	status := http.StatusOK
	generateFakeData(&fakeForm)
	fakeForm.Percent = nil
	fakeUserResponse := &models.Segment{
		Slug: fakeForm.Slug,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmentUC := mockSegmentUC.NewMockUseCaseI(ctrl)
	segmentH := New(cfg, segmentUC)

	body, err := json.Marshal(fakeForm)
	if err != nil {
		t.Fatalf("error while marshaling to json: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/segment/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	segmentUC.EXPECT().CreateSegment(fakeForm).Return(fakeUserResponse, nil)
	segmentH.CreateSegment(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected status %d, got %d ", status, w.Code)
	}
}

func TestDelivery_DeleteSegment(t *testing.T) {
	cfg := createConfig()

	slug := "test"
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmentUC := mockSegmentUC.NewMockUseCaseI(ctrl)
	segmentH := New(cfg, segmentUC)

	r := httptest.NewRequest(http.MethodDelete, "/segment/", bytes.NewReader([]byte{}))
	vars := map[string]string{
		"slug": slug,
	}

	r = mux.SetURLVars(r, vars)
	w := httptest.NewRecorder()

	segmentUC.EXPECT().DeleteSegment(slug).Return(nil)
	segmentH.DeleteSegment(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected status %d, got %d ", status, w.Code)
	}
}

func TestDelivery_GetSegment(t *testing.T) {
	cfg := createConfig()

	var fakeSegmentResponse *models.Segment
	generateFakeData(&fakeSegmentResponse)
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmentUC := mockSegmentUC.NewMockUseCaseI(ctrl)
	segmentH := New(cfg, segmentUC)

	r := httptest.NewRequest(http.MethodGet, "/segment/", bytes.NewReader([]byte{}))
	vars := map[string]string{
		"slug": fakeSegmentResponse.Slug,
	}

	r = mux.SetURLVars(r, vars)
	w := httptest.NewRecorder()

	segmentUC.EXPECT().GetSegmentBySlug(fakeSegmentResponse.Slug).Return(fakeSegmentResponse, nil)
	segmentH.GetSegment(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected status %d, got %d ", status, w.Code)
	}
}

func TestDelivery_GetUserSegments(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var fakeUserSegmentsResponse []models.Segment
	generateFakeData(&fakeUserSegmentsResponse)
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmentUC := mockSegmentUC.NewMockUseCaseI(ctrl)
	segmentH := New(cfg, segmentUC)

	r := httptest.NewRequest(http.MethodGet, "/user/{id}/segments", bytes.NewReader([]byte{}))
	vars := map[string]string{
		"id": strconv.FormatUint(userID, 10),
	}

	r = mux.SetURLVars(r, vars)
	w := httptest.NewRecorder()

	segmentUC.EXPECT().GetUserSegments(userID).Return(fakeUserSegmentsResponse, nil)
	segmentH.GetUserSegments(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected status %d, got %d ", status, w.Code)
	}
}

func TestDelivery_EditUserSegments(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var fakeForm models.FormEditSegments
	var fakeUserSegmentsResponse []models.Segment
	generateFakeData(&fakeForm)
	generateFakeData(&fakeUserSegmentsResponse)
	status := http.StatusOK

	for idx, _ := range fakeForm.SegmentsToAdd {
		fakeForm.SegmentsToAdd[idx].Until = nil
		fakeForm.SegmentsToAdd[idx].SegmentID = 0
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmentUC := mockSegmentUC.NewMockUseCaseI(ctrl)
	segmentH := New(cfg, segmentUC)

	body, err := json.Marshal(fakeForm)
	if err != nil {
		t.Fatalf("error while marshaling to json: %v", err)
	}

	r := httptest.NewRequest(http.MethodGet, "/user/{id}/segments", bytes.NewReader(body))
	vars := map[string]string{
		"id": strconv.FormatUint(userID, 10),
	}

	r = mux.SetURLVars(r, vars)
	w := httptest.NewRecorder()

	segmentUC.EXPECT().EditUserSegments(userID, fakeForm.SegmentsToAdd, fakeForm.SegmentsToRemove).Return(fakeUserSegmentsResponse, nil)
	segmentH.EditUserSegments(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected status %d, got %d ", status, w.Code)
	}
}
