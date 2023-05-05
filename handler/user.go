package handler

import (
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"jayant/database"
	"jayant/database/dbHelper"
	"jayant/middlewares"
	"jayant/models"
	"jayant/utils"
	"net/http"
)

// CreateUser handles the registration of a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	body := models.Register{}
	if Err := utils.ParseBody(r.Body, &body); Err != nil {
		utils.RespondError(w, http.StatusBadRequest, Err, "failed to parse request body")
		return
	}
	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "length of password is less than 6")
		return
	}
	exists, existsErr := dbHelper.IsEmailExits(body.Email)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check email exits")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "user already exists")
		return
	}
	hashedPassword, hasErr := utils.HashPassword(body.Password)
	if hasErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to create password")
		return
	}
	err := dbHelper.CreateUser(body, hashedPassword)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, "Registration of user is successful")
}

// LoginUser handles the login request of a user
func LoginUser(w http.ResponseWriter, r *http.Request) {
	body := models.Login{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	userid, err := dbHelper.RetrieveUserId(body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "username is wrong")
		return
	}
	err = dbHelper.MatchPassword(body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "password don't match")
		return
	}
	id := uuid.NewV1()
	sessionToken := id.String()
	sessionErr := dbHelper.CreateSession(userid, sessionToken)
	if sessionErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, sessionErr, "failed to create session")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"sessionToken": sessionToken,
	})
}

//UserInfo retrieves the information of the currently logged-in user
func UserInfo(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	UserInfo, err := dbHelper.UserInfo(user.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to retrieve User info")
		return
	}
	utils.RespondJSON(w, http.StatusOK, UserInfo)
}

//DeleteUser deletes the currently logged-in user and all associated sessions
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err := dbHelper.DeleteUser(user.ID)
		if err != nil {
			return err
		}
		err = dbHelper.DeleteSession(user.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to delete user")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UpdateUser updates the information of the authenticated user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	var body models.Register
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse request body")
		return
	}
	exists, existsErr := dbHelper.IsSameEmailUsedInOtherUser(body.Email, user.ID)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check email exits")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "email used by other user")
		return
	}
	err = database.Tx(func(tx *sqlx.Tx) error {
		err = dbHelper.UpdateUser(user.ID, body)
		if err != nil {
			return err
		}
		err = dbHelper.DeleteSession(user.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update user")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Logout deletes the session associated with the currently logged-in user
func Logout(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	err := dbHelper.DeleteSession(user.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to logout user")
		return
	}
	w.WriteHeader(http.StatusOK)
}
