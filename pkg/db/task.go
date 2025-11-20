package db

import (
	"database/sql"
	"errors"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Adding a new task into database
func AddTask(task *Task) (int64, error) {

	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	// Receiving the id of the inserted record
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

// Getting last <limit> tasks from a database
func Tasks(limit int, searchDate string, searchText string) ([]Task, error) {

	var tasks []Task
	var rows *sql.Rows
	var err error

	if searchDate != "" {

		// Executing search by date
		rows, err = DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date DESC LIMIT :limit",
			sql.Named("date", searchDate),
			sql.Named("limit", limit))

	}

	if searchText != "" {

		searchText = "%" + searchText + "%"

		// Executing search by text
		rows, err = DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :text OR comment LIKE :text ORDER BY date DESC LIMIT :limit",
			sql.Named("text", searchText),
			sql.Named("limit", limit))

	}

	if searchDate == "" && searchText == "" {

		// Search is not required, executing regular select statememt
		rows, err = DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date DESC LIMIT :limit",
			sql.Named("limit", limit))

	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {

		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)

	}

	// In case there are no records, returning {"tasks":[]} object
	if len(tasks) == 0 {
		return []Task{}, nil
	}

	return tasks, nil

}

// Searching for a single task record
func GetTask(id string) (Task, error) {

	var task Task
	var row *sql.Row

	row = DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	return task, err

}

// Updating of a single task record
func UpdateTask(task *Task) error {

	var (
		errIncorrectId = errors.New("Update failed: incorrect task id")
	)

	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`

	res, err := DB.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return err
	}

	// Checking the amount of affected rows
	count, err := res.RowsAffected()

	if err != nil {
		return err
	}

	// If nothing was updated, then incorrect task Id has been provided
	if count == 0 {
		return errIncorrectId
	}

	return nil

}

// Deletion of a single task record
func DeleteTask(id string) error {

	var (
		errIncorrectId = errors.New("Update failed: incorrect task id")
	)

	res, err := DB.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	if err != nil {
		return err
	}

	// Checking the amount of affected rows
	count, err := res.RowsAffected()

	if err != nil {
		return err
	}

	// If nothing was updated, then incorrect task Id has been provided
	if count == 0 {
		return errIncorrectId
	}

	return nil

}
