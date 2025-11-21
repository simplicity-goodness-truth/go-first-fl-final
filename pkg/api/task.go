package api

import (
	"net/http"
)

// Central task handler, which distributes functions depending on http method
func taskHandler(res http.ResponseWriter, req *http.Request) {

	switch req.Method {

	//  POST method	handler
	case http.MethodPost:
		addTaskHandler(res, req)

	//  GET method handler
	case http.MethodGet:
		getTaskHandler(res, req)

	//  PUT methodhandler
	case http.MethodPut:
		putTaskHandler(res, req)

	//  DELETE methodhandler
	case http.MethodDelete:
		deleteTaskHandler(res, req)

	default:
		http.Error(res, "Method is not permitted", http.StatusMethodNotAllowed)
		return
	}

}

// Central tasks handler, which distributes functions depending on http method

func tasksHandler(res http.ResponseWriter, req *http.Request) {

	switch req.Method {

	//  GET method handler
	case http.MethodGet:
		getTasksHandler(res, req)

	}

}
