package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"tracker/pkg/db"
)

// Task addition handler
func addTaskHandler(res http.ResponseWriter, req *http.Request) {

	// Deserializing the incoming json
	var task db.Task

	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		writeJson(res, "JSON decoding error", http.StatusBadRequest)
		return
	}

	// Validation of incoming title
	if task.Title == "" {
		writeJson(res, NewTaskResponse{Error: fmt.Sprintf("Empty title")}, http.StatusBadRequest)
		return
	}

	// Validation of incoming date
	err := checkDate(&task)

	if err != nil {
		writeJson(res, NewTaskResponse{Error: fmt.Sprintf("Date processing error: %v", err)}, http.StatusBadRequest)
		return
	}

	// Adding task to the database
	id, err := db.AddTask(&task)

	if err != nil {
		writeJson(res, NewTaskResponse{Error: fmt.Sprintf("Database insert error: %v", err)}, http.StatusInternalServerError)
		return
	}

	// Sending successfull response
	writeJson(res, NewTaskResponse{ID: fmt.Sprintf("%d", id)}, http.StatusOK)

}

// Checking date and setting task date
func checkDate(task *db.Task) error {

	var next string
	now := time.Now()

	// If task date is empty, then setting it to today
	if task.Date != "" {

		// Validating the correctness of a date format
		t, err := time.Parse(apiDateFormat, task.Date)

		if err != nil {
			return err
		}

		// Setting new dates
		nowDate := now.Format(apiDateFormat)
		taskDate := t.Format(apiDateFormat)

		// If task is later than now, we calculate the new date, otherwise setting it for today
		if task.Repeat != "" {

			next, err = NextDate(now, task.Date, task.Repeat)

			if err != nil {
				return err
			}

			if taskDate < nowDate   {
				task.Date = next
			} 	
		} else {
			if  taskDate < nowDate {
				task.Date = nowDate
			}
		}

	} else {
		task.Date = now.Format(apiDateFormat)
	}

	return nil

}
