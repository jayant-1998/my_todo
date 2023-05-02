package dbHelper

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"jayant/database"
	"jayant/models"
	"jayant/utils"
)

//func CreateUser(db sqlx.Ext, name, email, password string) (string, error) {
//	SQL := `insert into
//    				users(name,email,password)
//			values
//			    	($1,TRIM(LOWER($2)),$3)
//			returning
//						id`
//	var userID string
//
//	err := db.QueryRowx(SQL, name, email, password).Scan(&userID)
//	//fmt.Println(err)
//	if err != nil {
//		logrus.Errorf("CreateUser : error creating user: %v", err)
//		return "", err
//	}
//	return userID, nil
//}

func CreateUser(name, email, password string) error {
	SQL := `insert into 
					users(name,email,password)
			values 
				($1,trim(lower($2)),$3)`
	_, err := database.Todo.Exec(SQL, name, email, password)
	if err != nil {
		logrus.Errorf("CreateUser : error creating user: %v", err)
		return err
	}
	return nil
}

//EmailExits function returns bool
func EmailExits(email string) (bool, error) {
	SQL := `select id 
			from users 
			where 	email = $1 
			  		and archived_at is null`
	var id int
	err := database.Todo.Get(&id, SQL, email)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("EmailExits : error cheacking email exit or not: %v", err)
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

// InfoUser function returns user information
func InfoUser(id int) (models.UsersGetBody, error) {
	var body models.UsersGetBody
	SQL := `select  id,
       				name,
       				email,
       				created_at,
       				updated_at 
			from users
			where 	id = $1 
			  		and archived_at is null`
	err := database.Todo.Get(&body, SQL, id)
	//fmt.Println(err)
	if err != nil {
		logrus.Errorf("InfoUser : error  retrieving user info: %v", err)
		return body, err
	}
	return body, nil
}

func DeleteUser(id int) error {
	SQL := `update 
    				users 
			set 
			    	archived_at = now()
			where   id = $1 
			  		and archived_at is null `
	_, err := database.Todo.Exec(SQL, id)
	if err != nil {
		logrus.Errorf("DeleteUser : error deleting user: %v", err)
		return err
	}
	return nil
}

func DeleteUserSession(id int) error {
	SQL := `update
					user_session
			set
			    	archived_at = now()
			where
			    	user_id = $1`
	_, err := database.Todo.Exec(SQL, id)
	if err != nil {
		logrus.Errorf("Session_User : error session user: %v", err)
		return err
	}
	return nil
}

func UpdateUser(id int, name string, email string) error {

	SQL := `update users

			set 	name = $2,
			    	email=$3,
			    	updated_at = now()
			
			where 	id = $1 
			  		and archived_at is null `
	_, err := database.Todo.Exec(SQL, id, name, email)
	if err != nil {
		logrus.Errorf("UpdateUser : error updating user: %v", err)
		return err
	}
	return nil
}

func CreateUserSession(id int, SessionToken string) error {
	SQL := `insert into 
    					user_session(user_Id,session_token)
			values 
			    	($1,$2)`
	_, err := database.Todo.Exec(SQL, id, SessionToken)
	if err != nil {
		logrus.Errorf("CreateUserSession : error updating user session: %v", err)
		return err
	}
	return nil
}

/*
func UserIDByPassword(email, password string) (string, error) {
	s := `select
					id  as user_id,
					password
			  From
					users
			  where
				    archived_at is null
					and email = TRIM(LOWER($1))`

	//var userID string
	//var passwordHash string

	type UserPass struct {
		UserID   int    `db:"user_id"`
		Password string `db:"password"`
	}
	userDetails := UserPass{}

	err := database.Todo.Get(&userDetails, s, email)

	//fmt.Println(userDetails.UserID)
	//fmt.Println(Password)

	if err == sql.ErrNoRows {
		return "", nil
	}

	if passwordErr := utils.CheckPassword(password, userDetails.Password); passwordErr != nil {
		return "", passwordErr
	}
	return userID, nil
}
*/

func GetIDByEmail(email, password string) (int, error) {
	s := `select 
					id  as user_id,
					password
			  From
					users 
			  where
				    archived_at is null 
					and email = TRIM(LOWER($1))`

	type UserPass struct {
		UserID   int    `db:"user_id"`
		Password string `db:"password"`
	}
	userDetails := UserPass{}

	err := database.Todo.Get(&userDetails, s, email)

	//fmt.Println(userDetails.UserID)
	//fmt.Println(Password)

	if err == sql.ErrNoRows {
		return 0, nil
	}
	if passwordErr := utils.CheckPassword(password, userDetails.Password); passwordErr != nil {
		return 0, passwordErr
	}
	return userDetails.UserID, nil
}

func GetUserBySession(sessionToken string) (*models.User, error) {
	// language=SQL
	SQL := `SELECT 
       			u.id, 
       			u.name, 
       			u.email, 
       			u.created_at 
			FROM 
			    users u
			JOIN 
			    user_session us on u.id = us.user_id
			WHERE 	
			    u.archived_at IS NULL 
			  	AND us.session_token = $1
			  	AND us.archived_at > now()`
	var user models.User
	err := database.Todo.Get(&user, SQL, sessionToken)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, nil
}
