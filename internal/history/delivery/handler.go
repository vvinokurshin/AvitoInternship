package delivery

import (
	"github.com/go-playground/validator/v10"
	pkgErrors "github.com/pkg/errors"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	historyUC "github.com/vvinokurshin/AvitoInternship/internal/history/usecase"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	"github.com/vvinokurshin/AvitoInternship/pkg"
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

type DeliveryI interface {
	GetHistoryCSV(w http.ResponseWriter, r *http.Request)
}

type Delivery struct {
	cfg *config.Config
	uc  historyUC.UseCaseI
}

// GetHistoryCSV godoc
// @Summary      GetHistoryCSV
// @Description  getting the history of adding/removing users in a segment for a specific month
// @Tags     history
// @Accept	 application/json
// @Produce  application/json
// @Param year query int true "year"
// @Param month query int true "month"
// @Success 200 "downloaded file"
// @Failure 400 {object} errors.JSONError "year is required"
// @Failure 400 {object} errors.JSONError "year is invalid"
// @Failure 400 {object} errors.JSONError "month is required"
// @Failure 400 {object} errors.JSONError "month is invalid"
// @Failure 400 {object} errors.JSONError "invalid parameters"
// @Failure 500 {object} errors.JSONError "internal server error"
// @Router   /history [get]
func (d *Delivery) GetHistoryCSV(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		pkg.HandleError(w, r, errors.ErrYearIsRequired)
		return
	}
	month, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil {
		pkg.HandleError(w, r, errors.ErrMonthIsRequired)
		return
	}

	form := models.FormHistory{
		Year:  year,
		Month: time.Month(month),
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidParameters, err.Error()))
		return
	}

	err = form.Validate()
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	fileName, err := d.uc.GetHistoryCSV(form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	pkg.SendFile(w, r, fileName)
}

func New(cfg *config.Config, uc historyUC.UseCaseI) DeliveryI {
	return &Delivery{
		cfg: cfg,
		uc:  uc,
	}
}
