package api

import (
	"net/http"
	"strconv"

	"github.com/cruffinoni/neobrain-todolist/internal/database"
	"github.com/cruffinoni/neobrain-todolist/internal/utils"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	db *database.DB
}

func NewRoutes(db *database.DB) *Routes {
	return &Routes{db: db}
}

type AddTaskBody struct {
	Task string `json:"task"`
}

type AddTaskResponse struct {
	TaskID int64 `json:"task_id"`
}

func (r *Routes) AddTask(ctx *gin.Context) {
	var task AddTaskBody
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder(err.Error()))
		return
	}
	id, err := r.db.AddTask(task.Task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewInternalServerErrorBuilder(err))
		return
	}
	ctx.JSON(http.StatusCreated, AddTaskResponse{TaskID: id})
}

func (r *Routes) DeleteTask(ctx *gin.Context) {
	taskIDStr := ctx.Param("task_id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder("invalid task id"))
		return
	}

	err = r.db.DeleteTask(taskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewInternalServerErrorBuilder(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewStatusOKBuilder("task deleted successfully"))
}

func (r *Routes) MarkAsDone(ctx *gin.Context) {
	taskIDStr := ctx.Param("task_id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder("invalid task id"))
		return
	}

	err = r.db.MarkTaskAsDone(taskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewInternalServerErrorBuilder(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewStatusOKBuilder("task marked as done"))
}

func (r *Routes) ListTasks(ctx *gin.Context) {
	filterStr := ctx.Query("completed")

	var filter database.TaskFilter
	switch filterStr {
	case "true":
		filter = database.TaskFilterDone
	case "false":
		filter = database.TaskFilterNotDone
	default:
		filter = database.TaskFilterAll
	}

	tasks, err := r.db.GetTasks(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewInternalServerErrorBuilder(err))
		return
	}
	if len(tasks) == 0 {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
