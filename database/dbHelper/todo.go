package dbHelper

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"jayant/database"
	"jayant/models"
)

func CreateTask(id int, body models.Task) error {
	SQL := `INSERT INTO 
    					todo(user_id, task , description)
			VALUES 
			    	($1, $2, $3)`
	_, err := database.Todo.Exec(SQL, id, body.Task, body.Description)
	if err != nil {
		logrus.Errorf("CreateTodo: error creating task: %v", err)
		return err
	}
	return nil
}

func AllTask(id int) ([]models.AllTask, error) {
	body := make([]models.AllTask, 0)
	SQl := `SELECT 
    			id ,
				task,
				description,
				is_completed ,
				due_date ,
				created_at
			FROM 
			    todo
			WHERE 
			    user_id = $1
			  	AND archived_at IS NULL `
	err := database.Todo.Select(&body, SQl, id)
	if err != nil {
		logrus.Errorf("AllTask: failed to retrieve all task: %v", err)
		return nil, err
	}
	return body, nil
}

func UpdateTask(uid int, body models.UpdateTask) error {
	SQL := `UPDATE 
    				todo
			SET 
			    	description=$4,
			    	task = $3
			WHERE 
			    	id = $1 
			  		AND user_id = $2
			  		AND archived_at IS NULL`
	r, err := database.Todo.Exec(SQL, body.Id, uid, body.Task, body.Description)
	fmt.Println(r)
	if err != nil {
		logrus.Errorf("UpdateTask : failed to update task : %v", err)
		return err
	}
	return nil
}

func DeleteTask(id, uid int) error {
	SQL := `UPDATE 
    			todo
			SET 
			    archived_at = now()
			WHERE 
			    id = $1 
			  	AND user_id = $2
			  	AND archived_at IS NULL `
	_, err := database.Todo.Exec(SQL, id, uid)
	if err != nil {
		logrus.Errorf("DeleteTodo: failed to delete task: %v", err)
		return err
	}
	return nil
}

func Complete(id, uid int) error {
	SQL := `UPDATE 
    				todo 
			SET 
			    	is_completed = true 
			WHERE 
			    	id = $1
			  		AND user_id =$2
			    	AND is_completed = FALSE
			    	AND archived_at IS NULL `
	_, err := database.Todo.Exec(SQL, id, uid)
	if err != nil {
		logrus.Errorf("complete : Failed to update complete: %v", err)
		return err
	}
	return nil
}
