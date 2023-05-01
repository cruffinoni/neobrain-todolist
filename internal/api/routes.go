package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/cruffinoni/neobrain-todolist/internal/database"
	"github.com/cruffinoni/neobrain-todolist/internal/utils"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	db database.Database
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
	if task.Task == "" {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder("the task can't be empty"))
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

func (r *Routes) MarkTaskAsDone(ctx *gin.Context) {
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
	log.Printf("Ret from err: %v & %v", tasks, err)
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

func (r *Routes) ImportTasks(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder(err.Error()))
		return
	}
	defer file.Close()

	formats := map[string]Importer{
		"csv":  CSVFormat{},
		"xlsx": ExcelFormat{},
	}
	i, ok := formats[ctx.DefaultQuery("format", "none")]
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder("unknown format type"))
	}

	var tasks []*database.Task
	if tasks, err = i.ImportTasks(file); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder(err.Error()))
		return
	}

	err = r.db.ImportTasks(tasks)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewInternalServerErrorBuilder(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewStatusOKBuilder("tasks imported successfully"))
}

func (r *Routes) ExportTasks(ctx *gin.Context) {
	tasks, err := r.db.GetTasks(database.TaskFilterAll)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewInternalServerErrorBuilder(err))
		return
	}

	var e Exporter
	switch ctx.DefaultQuery("format", "none") {
	case "csv":
		e = CSVFormat{}
		ctx.Header("Content-Disposition", "attachment; filename=tasks.csv")
		ctx.Header("Content-Type", "text/csv")
	case "xlsx":
		e = ExcelFormat{}
		ctx.Header("Content-Disposition", "attachment; filename=tasks.xlsx")
		ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	default:
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder("unknown format type"))
		return
	}
	if err = e.ExportTasks(tasks, ctx.Writer); err != nil {
		return
	}
	ctx.Status(http.StatusOK)
}
