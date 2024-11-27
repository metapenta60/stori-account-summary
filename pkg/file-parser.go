package pkg

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"stori-account-summary/model"
	"strconv"
	"strings"
)

func parseFile(ioReader io.Reader) (model.Rows, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(ioReader)

	if err != nil {
		return model.Rows{}, err
	}

	content := buf.String()

	// Prepare return
	csvReader := csv.NewReader(strings.NewReader(content))
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Unable to parse file as CSV, %v", err)
		return model.Rows{}, err
	}

	rows := model.Rows{}
	for _, record := range records {
		dateArray := strings.Split(record[1], "/")
		value, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return model.Rows{}, err
		}

		row := model.Row{
			Transaction: value,
			Date: model.Date{
				Day:   dateArray[1],
				Month: dateArray[0],
			},
		}

		rows = append(rows, row)
	}

	return rows, nil
}
