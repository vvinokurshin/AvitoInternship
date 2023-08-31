package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	mockUserUC "github.com/vvinokurshin/AvitoInternship/internal/user/usecase/mocks"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type inputCase struct {
	userID   uint64
	userForm models.FormUser
}

type outputCase struct {
	status int
}

type testCases struct {
	name   string
	input  inputCase
	output outputCase
}

func createConfig() *config.Config {
	return new(config.Config)
}

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestDelivery_CreateUser(t *testing.T) {
	cfg := createConfig()

	var fakeForm models.FormUser
	status := http.StatusOK
	generateFakeData(&fakeForm)
	fakeUserResponse := &models.User{
		UserID:    1,
		Username:  fakeForm.Username,
		FirstName: fakeForm.FirstName,
		LastName:  fakeForm.LastName,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := New(cfg, userUC)

	body, err := json.Marshal(fakeForm)
	if err != nil {
		t.Fatalf("error while marshaling to json: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/user/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	userUC.EXPECT().CreateUser(fakeForm).Return(fakeUserResponse, nil)
	userH.CreateUser(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected status %d, got %d ", status, w.Code)
	}
}

func TestDelivery_EditUser(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var fakeForm models.FormUser
	status := http.StatusOK
	generateFakeData(&fakeForm)
	fakeUserResponse := &models.User{
		UserID:    userID,
		Username:  fakeForm.Username,
		FirstName: fakeForm.FirstName,
		LastName:  fakeForm.LastName,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := New(cfg, userUC)

	body, err := json.Marshal(fakeForm)
	if err != nil {
		t.Fatalf("error while marshaling to json: %v", err)
	}

	r := httptest.NewRequest(http.MethodPut, "/user/", bytes.NewReader(body))
	vars := map[string]string{
		"id": strconv.FormatUint(userID, 10),
	}

	r = mux.SetURLVars(r, vars)
	w := httptest.NewRecorder()

	userUC.EXPECT().EditUser(userID, fakeForm).Return(fakeUserResponse, nil)
	userH.EditUser(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected status %d, got %d ", status, w.Code)
	}
}

func TestDelivery_DeleteUser(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := New(cfg, userUC)

	r := httptest.NewRequest(http.MethodDelete, "/user/", bytes.NewReader([]byte{}))
	vars := map[string]string{
		"id": strconv.FormatUint(userID, 10),
	}

	r = mux.SetURLVars(r, vars)
	w := httptest.NewRecorder()

	userUC.EXPECT().DeleteUser(userID).Return(nil)
	userH.DeleteUser(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected status %d, got %d ", status, w.Code)
	}
}

func TestDelivery_GetUser(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var fakeUserResponse *models.User
	generateFakeData(&fakeUserResponse)
	fakeUserResponse.UserID = 1
	status := http.StatusOK

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mockUserUC.NewMockUseCaseI(ctrl)
	userH := New(cfg, userUC)

	r := httptest.NewRequest(http.MethodGet, "/user/", bytes.NewReader([]byte{}))
	vars := map[string]string{
		"id": strconv.FormatUint(userID, 10),
	}

	r = mux.SetURLVars(r, vars)
	w := httptest.NewRecorder()

	userUC.EXPECT().GetUserByID(userID).Return(fakeUserResponse, nil)
	userH.GetUser(w, r)

	if w.Code != status {
		t.Errorf("[TEST] simple: Expected status %d, got %d ", status, w.Code)
	}
}

//q := r.URL.Query()
//q.Set("fromFolder", test.input.fromFolder)
//r.URL.RawQuery = q.Encode()
