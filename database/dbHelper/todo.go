package dbHelper

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"jayant/database"
	"jayant/models"
)

func CreateTask(id int, task string, description string) error {
	SQL := `insert into 
    					todo(user_id, task , description)
			values 
			    	($1,$2,$3)`
	_, err := database.Todo.Exec(SQL, id, task, description)
	if err != nil {
		logrus.Errorf("CreateTodo: error creating task: %v", err)
		return err
	}
	return nil
}

func TodoInfo(id int) ([]models.TodoGetBody, error) {
	body := make([]models.TodoGetBody, 0)
	SQl := `select id,
				   user_id,
				   task,
				   description,
				   created_at,
				   due_date
			from todo
			where user_id = $1
			  and is_completed = false
			  and archived_at is null`
	err := database.Todo.Select(&body, SQl, id)
	//fmt.Println(err)
	if err != nil {
		logrus.Errorf("TodoInfo: error retrieving data : %v", err)
		return body, err
	}
	return body, nil
}

func UpdateTodo(id int, task string, des string) error {
	SQL := `update 
    				todo
			set 
			    	description=$3
			where 
			    	user_id = $1 
			  		and task =$2
			  		and archived_at is null
			  		and is_completed = false`
	_, err := database.Todo.Exec(SQL, id, task, des)
	if err != nil {
		logrus.Errorf("UpdateTodo : error updating task: %v", err)
		return err
	}
	return nil
}

func TaskExits(id int, task string) (bool, error) {
	SQL := `select id 
			from todo
			where 	user_id = $1 
			  		and task = $2
			  		and is_completed = false
			  		and archived_at is null`
	var i int
	err := database.Todo.Get(&i, SQL, id, task)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("EmailExits : error cheacking email exit or not: %v", err)
		return false, err
	}
	if err == sql.ErrNoRows {
		return true, nil
	}
	return false, nil
}

func DeleteTodo(name string, id int) error {
	SQL := `update todo
			set archived_at = current_timestamp
			where user_id = $2 and task = $1 and archived_at is null `
	_, err := database.Todo.Exec(SQL, name, id)
	if err != nil {
		logrus.Errorf("DeleteTodo : error deleting task: %v", err)
		return err
	}
	return nil
}

func IsComplete(id int, task string) error {
	SQL := `update 
    				todo 
			set 
			    	is_completed = true 
			where 
			    	user_id = $1
			    	and is_completed = false
			    	and task = $2
			    	and archived_at is null `
	_, err := database.Todo.Exec(SQL, id, task)
	if err != nil {
		logrus.Errorf("IsComplete : error updating completing task: %v", err)
		return err
	}
	return nil
}
