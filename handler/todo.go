package handler

import (
	"github.com/go-chi/chi/v5"
	"jayant/database/dbHelper"
	"jayant/middlewares"
	"jayant/models"
	"jayant/utils"
	"net/http"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	id := middlewares.UserContext(r)
	body := models.TodoBody{}
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse request body")
		return
	}
	err = dbHelper.CreateTask(id.ID, body.Task, body.Description)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to enter the data in todo table")
	}
	w.WriteHeader(http.StatusOK)
}

func GetTodoInfo(w http.ResponseWriter, r *http.Request) {
	// get path parameters
	//i := chi.URLParam(r, "id")
	//id, err := strconv.Atoi(i)
	//if err != nil {
	//	utils.RespondError(w, http.StatusInternalServerError, err, "Failed to convert your string to int")
	//}
	id := middlewares.UserContext(r)
	TodoInfo, err := dbHelper.TodoInfo(id.ID)
	if err != nil && TodoInfo == nil {
		utils.RespondError(w, http.StatusBadRequest, err, "this is can created any task")
		return
	} else if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to get your data")
		return
	} else {
		utils.RespondJSON(w, http.StatusOK, TodoInfo)
		return
	}
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// get path parameters
	//name := chi.URLParam(r, "name")
	//i := chi.URLParam(r, "id")
	//id, err := strconv.Atoi(i)
	//if err != nil {
	//	utils.RespondError(w, http.StatusInternalServerError, err, "Failed to convert your string to int")
	//	return
	//}
	id := middlewares.UserContext(r)
	var body models.TodoUpdateBody
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse request body")
		return
	}
	exists, existsErr := dbHelper.TaskExits(id.ID, body.Task)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, existsErr, "task does not exits")
		return
	}
	err = dbHelper.UpdateTodo(id.ID, body.Task, body.Description)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update todo")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func IsComplete(w http.ResponseWriter, r *http.Request) {
	task := chi.URLParam(r, "task")
	//id, err := strconv.Atoi(i)
	//if err != nil {
	//	utils.RespondError(w, http.StatusBadRequest, err, "user enter string in place of integer")
	//	return
	//}
	id := middlewares.UserContext(r)

	err := dbHelper.IsComplete(id.ID, task)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update todo")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// get path parameters
	name := chi.URLParam(r, "name")
	//i := chi.URLParam(r, "id")
	//id, err := strconv.Atoi(i)
	//if err != nil {
	//	utils.RespondError(w, http.StatusBadRequest, err, "user enter string in place of integer")
	//	return
	//}
	id := middlewares.UserContext(r)
	Err := dbHelper.DeleteTodo(name, id.ID)
	if Err != nil {
		utils.RespondError(w, http.StatusInternalServerError, Err, "failed to delete todo")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
