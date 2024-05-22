package routes

import (
	"github.com/badaccuracyid/softeng_backend/src/controllers"
	"github.com/badaccuracyid/softeng_backend/src/database"
	"github.com/badaccuracyid/softeng_backend/src/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRoutes struct {
	baseRouter     *gin.RouterGroup
	userController controllers.UserController
}

func NewUserRoutes(router *gin.Engine) (*UserRoutes, error) {
	postgresDatabase, err := database.GetPostgresDatabase()
	if err != nil {
		panic(err)
	}

	userService := controllers.NewUserService(postgresDatabase)
	baseRouter := router.Group("/api/v1/users")

	return &UserRoutes{
		baseRouter:     baseRouter,
		userController: userService,
	}, nil
}

func (u *UserRoutes) InitializeRoutes() {
	u.baseRouter.GET("/", u.getUsers)
	u.baseRouter.POST("/", u.createUser)
}

// getUsers handles the GET /api/v1/users request
// @Summary Get all users
// @Description Get details of all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} model.User
// @Failure 500 {string} string
// @Router /api/v1/users [get]
func (u *UserRoutes) getUsers(ctx *gin.Context) {
	u.userController.SetContext(ctx)
	users, err := u.userController.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// createUser handles the POST /api/v1/users request
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body model.User true "User"
// @Success 201 {object} model.User
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/v1/users [post]
func (u *UserRoutes) createUser(ctx *gin.Context) {
	u.userController.SetContext(ctx)
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	createdUser, err := u.userController.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, createdUser)
}

func InitializeUserRoutes(router *gin.Engine) *gin.Engine {
	userRoutes, err := NewUserRoutes(router)
	if err != nil {
		panic(err)
	}

	userRoutes.InitializeRoutes()
	return router
}
