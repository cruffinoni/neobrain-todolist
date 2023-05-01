package api

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

	"github.com/cruffinoni/neobrain-todolist/internal/database"
	"github.com/xuri/excelize/v2"
)

type Importer interface {
	ImportTasks(f multipart.File) ([]*database.Task, error)
}

type Exporter interface {
	ExportTasks(tasks []*database.Task, w io.Writer) error
}

func extractTaskFromRow(row []string) (*database.Task, error) {
	if len(row) != 2 {
		return nil, fmt.Errorf("malformed row: expected only 2 arg, got: %d", len(row))
	}
	var (
		err error
		t   = &database.Task{
			Task: row[0],
		}
	)
	t.Done, err = strconv.ParseBool(row[1])
	if err != nil {
		return nil, errors.New("invalid field in 'done' section. Value must either true or false, got '" + row[1] + "'")
	}
	return t, nil
}

type CSVFormat struct{}

func (CSVFormat) ImportTasks(f multipart.File) ([]*database.Task, error) {
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks []*database.Task
	for _, record := range records[1:] {
		if t, err := extractTaskFromRow(record); err != nil {
			return nil, err
		} else {
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

func (CSVFormat) ExportTasks(tasks []*database.Task, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	err := writer.Write([]string{"task", "done"})
	if err != nil {
		return err
	}
	for _, task := range tasks {
		err := writer.Write([]string{task.Task, strconv.FormatBool(task.Done)})
		if err != nil {
			return err
		}
	}

	return nil
}

type ExcelFormat struct{}

func (ExcelFormat) ImportTasks(f multipart.File) ([]*database.Task, error) {
	wb, err := excelize.OpenReader(f)
	if err != nil {
		return nil, err
	}

	rows, err := wb.GetRows(wb.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	var tasks []*database.Task
	for _, row := range rows[1:] {
		if t, err := extractTaskFromRow(row); err != nil {
			return nil, err
		} else {
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}

func (ExcelFormat) applyValueToCell(sheetName string, file *excelize.File, col, row int, value any) error {
	cell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, cell, value)
	if err != nil {
		return err
	}
	return nil
}

func (e ExcelFormat) ExportTasks(tasks []*database.Task, w io.Writer) error {
	f := excelize.NewFile()
	sheetName := "Sheet1"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	if err := e.applyValueToCell(sheetName, f, 1, 1, "task"); err != nil {
		return err
	}
	if err := e.applyValueToCell(sheetName, f, 2, 1, "done"); err != nil {
		return err
	}

	for i, task := range tasks {
		if err := e.applyValueToCell(sheetName, f, 1, i+2, task.Task); err != nil {
			return err
		}
		if err := e.applyValueToCell(sheetName, f, 2, i+2, task.Done); err != nil {
			return err
		}
	}

	return f.Write(w)
}
