package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"tracker/pkg/db"
)

// Handler to get a tasks list
func getTasksHandler(res http.ResponseWriter, req *http.Request) {

	var tasks []db.Task
	var err error

	// Getting search parameter
	search := req.FormValue("search")

	if search != "" {

		// Checking if it search by date
		searchDate, err := time.Parse("02.01.2006", search)

		if err == nil {

			// Switching date to a unified api format and
			// fetching last 50 tasks from a database by date
			tasks, err = db.Tasks(50, searchDate.Format(apiDateFormat), "")

		} else {

			// Fetching last 50 tasks from a database by text
			tasks, err = db.Tasks(50, "", search)
		}

	} else {

		// Fetching last 50 tasks from a database
		tasks, err = db.Tasks(50, "", "")

	}

	if err != nil {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("database selection failed")}, http.StatusInternalServerError)
		return
	}

	// Sending a successfull response to a client
	writeJson(res, TasksResponse{Tasks: tasks}, http.StatusOK)

}

// Handler to update a task list through a PUT method
func putTaskHandler(res http.ResponseWriter, req *http.Request) {

	// Deserializing the incoming json
	var task db.Task

	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		writeJson(res, "JSON decoding error", http.StatusBadRequest)
		return
	}

	// Validation of incoming title
	if task.Title == "" {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("Empty title")}, http.StatusBadRequest)
		return
	}

	// Validation of incoming date
	err := checkDate(&task)

	if err != nil {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("Date processing error: %v", err)}, http.StatusBadRequest)
		return
	}

	// Updating a task in the database
	err = db.UpdateTask(&task)

	if err != nil {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("Database update error: %v", err)}, http.StatusInternalServerError)
		return
	}

	// Sending successfull response
	writeJson(res, EmptyResponse{}, http.StatusOK)

}

// Handler to get a single tasks list
func getTaskHandler(res http.ResponseWriter, req *http.Request) {

	var task db.Task
	var err error

	// Getting ID parameter
	id := req.FormValue("id")

	//ID parameter cannot be empty
	if id == "" {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("task id has not been provided")}, http.StatusBadRequest)
		return
	}

	// Fetching a task by ID from a database
	task, err = db.GetTask(id)

	if err != nil {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("database selection failed")}, http.StatusInternalServerError)
		return
	}

	// Sending a response in case a task has not been found
	if (task == db.Task{}) {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("task has not been found")}, http.StatusBadRequest)
		return
	}

	// Sending a successfull response to a client
	writeJson(res, task, http.StatusOK)
}

// Handler to update a task list through a DELETE method
func deleteTaskHandler(res http.ResponseWriter, req *http.Request) {

	var err error

	// Getting ID parameter
	id := req.FormValue("id")

	//ID parameter cannot be empty
	if id == "" {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("task id has not been provided")}, http.StatusBadRequest)
		return
	}

	// Deleting a task by ID from a database
	err = db.DeleteTask(id)

	if err != nil {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("database selection failed")}, http.StatusInternalServerError)
		return
	}

	// Sending successfull response
	writeJson(res, EmptyResponse{}, http.StatusOK)

}

// Handler to mark a task as done
func taskDoneHandler(res http.ResponseWriter, req *http.Request) {

	var task db.Task
	var err error

	// Getting ID parameter
	id := req.FormValue("id")

	//ID parameter cannot be empty
	if id == "" {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("task id has not been provided")}, http.StatusBadRequest)
		return
	}

	// Getting task by ID from the database

	// Fetching a task by ID from a database
	task, err = db.GetTask(id)

	if err != nil {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("database selection failed")}, http.StatusInternalServerError)
		return
	}

	// Sending a response in case a task has not been found
	if (task == db.Task{}) {
		writeJson(res, ErrorResponse{Error: fmt.Sprintf("task has not been found")}, http.StatusBadRequest)
		return
	}

	// Validation if task is repeatable

	if task.Repeat != "" {

		now := time.Now()

		// Calculating a next date
		next, err := NextDate(now, task.Date, task.Repeat)

		if err != nil {
			writeJson(res, ErrorResponse{Error: fmt.Sprintf("Next date calculation error: %v", err)}, http.StatusInternalServerError)
			return
		}

		// Updating a task in the database

		task.Date = next
		err = db.UpdateTask(&task)

		// Updating a task in the database
		err = db.UpdateTask(&task)

		if err != nil {
			writeJson(res, ErrorResponse{Error: fmt.Sprintf("Database update error: %v", err)}, http.StatusInternalServerError)
			return
		}

		// Sending successfull response
		writeJson(res, EmptyResponse{}, http.StatusOK)

	} else {
		// If the task not repeatable, then we are deleting it

		err = db.DeleteTask(id)

		if err != nil {
			writeJson(res, ErrorResponse{Error: fmt.Sprintf("database selection failed")}, http.StatusInternalServerError)
			return
		}

		// Sending successfull response
		writeJson(res, EmptyResponse{}, http.StatusOK)

	}

}
