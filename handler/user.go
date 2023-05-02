package handler

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"jayant/database"
	"jayant/database/dbHelper"
	"jayant/middlewares"
	"jayant/models"
	"jayant/utils"
	"net/http"
	"time"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body := models.UsersBody{}
	if Err := utils.ParseBody(r.Body, &body); Err != nil {
		utils.RespondError(w, http.StatusBadRequest, Err, "failed to parse request body")
		return
	}
	// check  the length of password
	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "length of password is less than 6")
	}
	// check the email is already  exits or not
	exists, existsErr := dbHelper.EmailExits(body.Email)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "user already exists")
		return
	}
	// convert password into hash password
	hashedPassword, hasErr := utils.HashPassword(body.Password)
	if hasErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to create user")
		return
	}
	err := dbHelper.CreateUser(body.Name, body.Email, hashedPassword)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user")
		return
	}
	// creating a session for user
	//sessionToken := utils.HashString(body.Email + time.Now().String())
	//txErr := database.Tx(func(tx *sqlx.Tx) error {
	//	_, saveErr := dbHelper.CreateUser(tx, body.Name, body.Email, hashedPassword)
	//	if saveErr != nil {
	//		return saveErr
	//	}
	//	/*
	//		sessionErr := dbHelper.CreateUserSession(tx, userID, sessionToken)
	//		if sessionErr != nil {
	//			return sessionErr
	//		}
	//	*/
	//	return nil
	//})
	//if txErr != nil {
	//	utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create user")
	//	return
	//}
	//this is showing session
	utils.RespondJSON(w, http.StatusCreated, "Registration of user is successful")

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	userID, userErr := dbHelper.GetIDByEmail(body.Email, body.Password)
	if userErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, userErr, "failed to create user session")
		return
	}
	// converting  string to int id
	/*userID, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to convert given from UserIDByPassword")
		return
	}*/

	sessionToken := utils.HashString(body.Email + time.Now().String())
	sessionErr := dbHelper.CreateUserSession(userID, sessionToken)
	if sessionErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, sessionErr, "failed to create user session")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{
		Token: sessionToken,
	})
}

func InfoUser(w http.ResponseWriter, r *http.Request) {
	// get query parameters
	//q := r.URL.Query().Get("id")
	// get path parameters
	/*i := chi.URLParam(r, "id")
	id, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to convert your string to int")
		return
	}*/
	id := middlewares.UserContext(r)
	UserInfo, err := dbHelper.InfoUser(id.ID)
	// if no user id found for user
	if err == sql.ErrNoRows {
		utils.RespondError(w, http.StatusBadRequest, err, "User Enter Wrong ID")
		return
	} else if err != nil { // if any error occurs
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to get the User data")
		return
	}
	utils.RespondJSON(w, http.StatusOK, UserInfo)

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// get path parameters
	//i := chi.URLParam(r, "id")
	//id, err := strconv.Atoi(i)
	//if err != nil {
	//	utils.RespondError(w, http.StatusInternalServerError, err, "Failed to convert your string to int")
	//	return
	//}

	id := middlewares.UserContext(r)
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err := dbHelper.DeleteUser(id.ID)
		if err != nil {
			return err
		}
		err = dbHelper.DeleteUserSession(id.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to delete user")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// get path parameters
	//i := chi.URLParam(r, "id")
	//id, err := strconv.Atoi(i)
	//if err != nil {
	//	utils.RespondError(w, http.StatusInternalServerError, err, "Failed to convert your string to int")
	//	return
	//}
	id := middlewares.UserContext(r)
	var body models.UsersBody
	if Err := utils.ParseBody(r.Body, &body); Err != nil {
		utils.RespondError(w, http.StatusBadRequest, Err, "failed to parse request body")
		return
	}
	exists, existsErr := dbHelper.EmailExits(body.Email)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "email is already exists")
		return
	}
	Err := dbHelper.UpdateUser(id.ID, body.Name, body.Email)
	if Err != nil {
		utils.RespondError(w, http.StatusInternalServerError, Err, "failed to update user")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	id := middlewares.UserContext(r)
	err := dbHelper.DeleteUserSession(id.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to logout")
		return
	}
	w.WriteHeader(http.StatusOK)
}
