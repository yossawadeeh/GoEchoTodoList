package main

import (
	"echo-todolist/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataHandler struct {
	DB *gorm.DB
}

func (h *DataHandler) Initialize() {
	db, err := gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/todolist?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	h.DB = db
}

func (h *DataHandler) HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func (h *DataHandler) GetAllToDoList(c echo.Context) error {
	userId := c.Param("id")
	todolist := []model.ToDoList{}

	if err := h.DB.First(&todolist, userId).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  http.StatusNotFound,
			"message": "Not Found",
		})
	}

	return c.JSON(http.StatusOK, todolist)
}

func (h *DataHandler) GetToDoList(c echo.Context) error {
	userId := c.Param("id")
	uId, _ := strconv.ParseUint(userId, 10, 32)
	todoId := c.Param("todoId")
	tId, _ := strconv.ParseUint(todoId, 10, 32)
	todolist := model.ToDoList{}

	if err := h.DB.Joins("inner join users u on to_do_lists.user_id = u.id").
		Where("to_do_lists.id = ? and u.id = ?", tId, uId).
		First(&todolist).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  http.StatusNotFound,
			"message": "Not Found",
		})
	}

	return c.JSON(http.StatusOK, todolist)
}

func (h *DataHandler) CreateToDoList(c echo.Context) error {
	todo := model.ToDoList{}
	user := model.User{}
	status := model.Status{}

	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  http.StatusBadRequest,
			"message": "Bad Request",
		})
	}

	userErr := h.DB.First(&user, todo.UserId).Error
	statusErr := h.DB.First(&status, todo.StatusId).Error
	if userErr != nil || statusErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  http.StatusBadRequest,
			"message": "Bad Request",
		})
	}

	todo.CreateDate = time.Now()
	if err := h.DB.Save(&todo).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, todo)
}

func (h *DataHandler) UpdateToDoList(c echo.Context) error {
	todo := model.ToDoList{}
	todoId := c.Param("todoId")
	tId, _ := strconv.ParseUint(todoId, 10, 32)
	todoUpdate := model.ToDoList{}

	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  http.StatusBadRequest,
			"message": "Bad Request",
		})
	}

	fmt.Println("todoId", todoId)
	fmt.Println("UserId", todo.UserId)
	if err := h.DB.Where("id = ? AND user_id = ?", todoId, todo.UserId).First(&todoUpdate).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  http.StatusNotFound,
			"message": "Not Found",
		})
	}

	todoUpdate.ID = uint(tId)
	todoUpdate.Discription = todo.Discription
	todoUpdate.StatusId = todo.StatusId
	todoUpdate.CreateDate = todo.CreateDate
	todoUpdate.UserId = todo.UserId

	if err := h.DB.Save(&todoUpdate).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, todo)
}

func (h *DataHandler) DeleteToDoList(c echo.Context) error {
	userId := c.Param("id")
	uId, _ := strconv.ParseUint(userId, 10, 32)
	todoId := c.Param("todoId")
	tId, _ := strconv.ParseUint(todoId, 10, 32)
	todolist := model.ToDoList{}

	if err := h.DB.Joins("inner join users u on to_do_lists.user_id = u.id").
		Where("to_do_lists.id = ? and u.id = ?", tId, uId).
		First(&todolist).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  http.StatusNotFound,
			"message": "Not Found",
		})
	}

	if err := h.DB.Delete(&todolist).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New() // เพื่อสร้าง instance ของ Echo ขึ้นมาใช้งาน

	h := DataHandler{}
	h.Initialize()

	e.GET("/helloworld", h.HelloWorld)
	e.GET("/all-todolist/:id", h.GetAllToDoList)
	e.GET("/todolist/:id/:todoId", h.GetToDoList)
	e.POST("/todolist", h.CreateToDoList)
	e.PUT("/todolist/:todoId", h.UpdateToDoList)
	e.DELETE("/todolist/:id/:todoId", h.DeleteToDoList)

	e.Logger.Fatal(e.Start(":8080"))
}
