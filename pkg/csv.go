package pkg

import (
	"encoding/csv"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	"os"
	"strconv"
)

func WriteCSV(records []models.History, fileName string) error {
	file, err := os.Create(HistoryFolderName + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ';'

	err = writer.Write([]string{"user_id", "slug", "operation", "datetime"})
	if err != nil {
		return err
	}

	// Записываем данные структур в CSV
	for _, record := range records {
		err = writer.Write([]string{strconv.FormatUint(record.UserID, 10), record.SegmentSlug, record.Operation, record.Datetime})
		if err != nil {
			return err
		}
	}

	return nil
}
