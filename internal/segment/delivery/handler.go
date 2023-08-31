package delivery

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	pkgErrors "github.com/pkg/errors"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	segmentUC "github.com/vvinokurshin/AvitoInternship/internal/segment/usecase"
	"github.com/vvinokurshin/AvitoInternship/pkg"
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
	"net/http"
	"strconv"
)

type DeliveryI interface {
	CreateSegment(w http.ResponseWriter, r *http.Request)
	DeleteSegment(w http.ResponseWriter, r *http.Request)
	GetSegment(w http.ResponseWriter, r *http.Request)
	GetUserSegments(w http.ResponseWriter, r *http.Request)
	EditUserSegments(w http.ResponseWriter, r *http.Request)
}

type Delivery struct {
	cfg *config.Config
	uc  segmentUC.UseCaseI
}

func New(cfg *config.Config, uc segmentUC.UseCaseI) DeliveryI {
	return &Delivery{
		cfg: cfg,
		uc:  uc,
	}
}

// CreateSegment godoc
// @Summary      CreateSegment
// @Description  create segment
// @Tags     segment
// @Accept	 application/json
// @Produce  application/json
// @Param    segment body models.FormSegment true "form segment"
// @Success 200 {object} models.SegmentResponse "segment created"
// @Failure 400 {object} errors.JSONError "invalid form"
// @Failure 400 {object} errors.JSONError "percent is invalid"
// @Failure 409 {object} errors.JSONError "segment with this slug already exists"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /segment/create [post]
func (d *Delivery) CreateSegment(w http.ResponseWriter, r *http.Request) {
	form := models.FormSegment{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	err := form.Validate()
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	response, err := d.uc.CreateSegment(form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.SegmentResponse{
		Segment: *response,
	})
}

// DeleteSegment godoc
// @Summary      DeleteSegment
// @Description  delete segment
// @Tags     segment
// @Accept	 application/json
// @Produce  application/json
// @Param slug path string true "slug"
// @Success 200 "segment deleted"
// @Failure 400 {object} errors.JSONError "invalid url"
// @Failure 404 {object} errors.JSONError "segment not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /segment/{slug} [delete]
func (d *Delivery) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok {
		pkg.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	err := d.uc.DeleteSegment(slug)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetSegment godoc
// @Summary      GetSegment
// @Description  get segment
// @Tags     segment
// @Accept	 application/json
// @Produce  application/json
// @Param slug path string true "slug"
// @Success 200 {object} models.SegmentResponse "success get segment info"
// @Failure 400 {object} errors.JSONError "invalid url"
// @Failure 404 {object} errors.JSONError "segment not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /segment/{slug} [get]
func (d *Delivery) GetSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok {
		pkg.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	response, err := d.uc.GetSegmentBySlug(slug)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.SegmentResponse{
		Segment: *response,
	})
}

// GetUserSegments godoc
// @Summary      GetUserSegments
// @Description  get user's segment
// @Tags     segment
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "id"
// @Success 200 {object} models.SegmentsResponse "success get user's segments"
// @Failure 400 {object} errors.JSONError "invalid url"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/{id}/segments [get]
func (d *Delivery) GetUserSegments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkg.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	segments, err := d.uc.GetUserSegments(userID)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.SegmentsResponse{
		Segments: segments,
		Count:    len(segments),
	})
}

// EditUserSegments godoc
// @Summary      EditUserSegments
// @Description  edit user's segment
// @Tags     segment
// @Accept	 application/json
// @Produce  application/json
// @Param id path int true "id"
// @Param    segment body models.FormEditSegments true "form segment"
// @Success 200 {object} models.SegmentsResponse "success edit user's segments"
// @Failure 400 {object} errors.JSONError "invalid url"
// @Failure 400 {object} errors.JSONError  "field until is invalid. format: YYYY-MM-DD HH:MM"
// @Failure 404 {object} errors.JSONError "user not found"
// @Failure 404 {object} errors.JSONError "segment not found"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /user/{id}/segments/edit [put]
func (d *Delivery) EditUserSegments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		pkg.HandleError(w, r, errors.ErrInvalidURL)
		return
	}

	form := models.FormEditSegments{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	err = form.Validate()
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	segments, err := d.uc.EditUserSegments(userID, form.SegmentsToAdd, form.SegmentsToRemove)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	pkg.SendJSON(w, r, http.StatusOK, models.SegmentsResponse{
		Segments: segments,
		Count:    len(segments),
	})
}
