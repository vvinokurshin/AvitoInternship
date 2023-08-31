package usecase

import (
	pkgErr "github.com/pkg/errors"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	historyRepository "github.com/vvinokurshin/AvitoInternship/internal/history/repository"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	"github.com/vvinokurshin/AvitoInternship/pkg"
	"strconv"
)

//go:generate mockgen -destination=./mocks/usecase.go -source=./usecase.go -package=mocks

type UseCaseI interface {
	GetHistoryCSV(form models.FormHistory) (string, error)
}

type UseCase struct {
	cfg         *config.Config
	historyRepo historyRepository.RepositoryI
}

func New(cfg *config.Config, historyRepo historyRepository.RepositoryI) UseCaseI {
	return &UseCase{
		cfg:         cfg,
		historyRepo: historyRepo,
	}
}

func (uc *UseCase) GetHistoryCSV(form models.FormHistory) (string, error) {
	records, err := uc.historyRepo.SelectRecordsByDate(form.Year, form.Month)
	if err != nil {
		return "", pkgErr.Wrap(err, "delete segments from user")
	}

	fileName := "history-" + strconv.Itoa(form.Year) + "-" + strconv.Itoa(int(form.Month)) + ".csv"
	err = pkg.WriteCSV(records, fileName)
	if err != nil {
		return "", pkgErr.Wrap(err, "write csv")
	}

	return fileName, nil
}
