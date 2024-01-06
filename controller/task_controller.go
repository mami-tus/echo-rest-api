package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type ITaskController interface {
	GetAllTasks(c echo.Context) error
	GetTaskById(c echo.Context) error
	CreateTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	DeleteTask(c echo.Context) error
}

type taskController struct {
	tu usecase.ITaskUsecase
}

func NewTaskController(tu usecase.ITaskUsecase) ITaskController {
	return &taskController{tu}
}

func (tc *taskController) GetAllTasks(c echo.Context) error {
	/*
		HTTPリクエストからuserというキーに紐づいたデータを取得します。このデータはJWT（JSON Web Token）のトークンです。
		.(*jwt.Token)は型アサーションで、取得したデータを*jwt.Token（JWTトークンへのポインタ）として扱います。
		※ 型アサーション（.(型)）は、interface{}型から具体的な型への変換を試みる操作
	*/
	user := c.Get("user").(*jwt.Token)
	/*
		JWTトークンからクレームを取得します。クレームはトークンに含まれるユーザー情報です。
		.(jwt.MapClaims)は型アサーションで、クレームをjwt.MapClaims型（マップ型）として扱います。
	*/
	claims := user.Claims.(jwt.MapClaims)
	// claimsマップから"user_id"キーに紐づいた値を取得し、userIdに格納します。
	userId := claims["user_id"]

	taskRes, err := tc.tu.GetAllTasks(uint(userId.(float64)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskRes)
}

func (tc *taskController) GetTaskById(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	// HTTPリクエストのパラメータから"taskId"に該当する値を取得します。
	id := c.Param("taskId")
	// 取得したタスクID（文字列）を整数に変換します。strconv.Atoiは文字列を整数（int型）に変換する関数です。
	taskId, _ := strconv.Atoi(id)
	taskRes, err := tc.tu.GetTaskById(uint(userId.(float64)), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskRes)
}

func (tc *taskController) CreateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	task := model.Task{}
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	task.UserId = uint(userId.(float64))
	taskRes, err := tc.tu.CreateTask(task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, taskRes)
}

func (tc *taskController) UpdateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	// HTTPリクエストのパラメータから"taskId"に該当する値を取得します。
	id := c.Param("taskId")
	// 取得したタスクID（文字列）を整数に変換します。strconv.Atoiは文字列を整数（int型）に変換する関数です。
	taskId, _ := strconv.Atoi(id)
	// model.Task型の新しいインスタンスtaskを作成
	task := model.Task{}
	// c.Bind(&task)を使用して、リクエストボディからのデータをtask変数にマッピング
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	taskRes, err := tc.tu.UpdateTask(task, uint(userId.(float64)), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskRes)
}

func (tc *taskController) DeleteTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	// HTTPリクエストのパラメータから"taskId"に該当する値を取得します。
	id := c.Param("taskId")
	// 取得したタスクID（文字列）を整数に変換します。strconv.Atoiは文字列を整数（int型）に変換する関数です。
	taskId, _ := strconv.Atoi(id)
	err := tc.tu.DeleteTask(uint(userId.(float64)), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
