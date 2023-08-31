package delivery

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	pkgErrors "github.com/pkg/errors"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	userUC "github.com/vvinokurshin/AvitoInternship/internal/user/usecase"
	"github.com/vvinokurshin/AvitoInternship/pkg"
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
	"net/http"
	"strconv"
)

type DeliveryI interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	EditUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
}

type Delivery struct {
	cfg *config.Config
	uc  userUC.UseCaseI
}

func New(cfg *config.Config, uc userUC.UseCaseI) DeliveryI {
	return &Delivery{
		cfg: cfg,
		uc:  uc,
	}
}

// CreateUser godoc
// @Summary      CreateUser
// @Description  create user
// @Tags     user
// @Accept	 application/json
// @Produce  application/json
// @Param    segment body models.FormUser true "form user"
// @Success 200 {object} models.UserResponse "user created"
// @Failure 400 {object} errors.JSONError "invalid form"
// @Failure 409 {object} errors.JSONError "user with this nickname already exists"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/create [post]
func (d *Delivery) CreateUser(w http.ResponseWriter, r *http.Request) {
	form := models.FormUser{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	response, err := d.uc.CreateUser(form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.UserResponse{
		User: *response,
	})
}

// EditUser godoc
// @Summary      EditUser
// @Description  edit user
// @Tags     user
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "id"
// @Param    segment body models.FormUser true "form user"
// @Success 200 {object} models.UserResponse "user info edited"
// @Failure 400 {object} errors.JSONError "invalid url"
// @Failure 400 {object} errors.JSONError "invalid form"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 409 {object} errors.JSONError "user with this nickname already exists"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/{id} [put]
func (d *Delivery) EditUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkg.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	form := models.FormUser{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	response, err := d.uc.EditUser(userID, form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.UserResponse{
		User: *response,
	})
}

// DeleteUser godoc
// @Summary      DeleteUser
// @Description  delete user
// @Tags     user
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "id"
// @Success 200 "user deleted"
// @Failure 400 {object} errors.JSONError "invalid url"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/{id} [delete]
func (d *Delivery) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkg.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	err = d.uc.DeleteUser(userID)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetUser godoc
// @Summary      GetUser
// @Description  get user
// @Tags     user
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "id"
// @Success 200 {object} models.UserResponse "success get user info"
// @Failure 400 {object} errors.JSONError "invalid url"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/{id} [get]
func (d *Delivery) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkg.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	response, err := d.uc.GetUserByID(userID)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.UserResponse{
		User: *response,
	})
}
