package handler

import (
	"github.com/go-chi/chi/v5"
	"jayant/database/dbHelper"
	"jayant/middlewares"
	"jayant/models"
	"jayant/utils"
	"net/http"
	"strconv"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	body := models.Task{}
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse request body")
		return
	}
	err = dbHelper.CreateTask(user.ID, body)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create task")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AllTask(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	task, err := dbHelper.AllTask(user.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to retrieve all task")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"task": task,
	})
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	var body models.UpdateTask
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse request body")
		return
	}

	err = dbHelper.UpdateTask(user.ID, body)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update todo in database")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Complete(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	i := chi.URLParam(r, "taskId")
	id, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to convert string to int")
		return
	}
	err = dbHelper.Complete(id, user.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update complete")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	i := chi.URLParam(r, "taskId")
	id, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to convert string to int")
		return
	}
	Err := dbHelper.DeleteTask(id, user.ID)
	if Err != nil {
		utils.RespondError(w, http.StatusInternalServerError, Err, "failed to delete task")
		return
	}
	w.WriteHeader(http.StatusOK)
}
