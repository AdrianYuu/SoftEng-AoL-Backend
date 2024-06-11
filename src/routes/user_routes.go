package routes

import (
	"net/http"

	"github.com/badaccuracyid/softeng_backend/src/controllers"
	"github.com/badaccuracyid/softeng_backend/src/database"
	"github.com/badaccuracyid/softeng_backend/src/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func InitializeUserRoutes(router *gin.Engine) *gin.Engine {
	userRoutes, err := NewUserRoutes(router)
	if err != nil {
		panic(err)
	}

	userRoutes.registerRoutes()
	return router
}

func (u *UserRoutes) registerRoutes() {
	u.baseRouter.POST("/create", u.createUser)
	u.baseRouter.PATCH("/update", u.createUser)

	u.baseRouter.GET("/get/:id", u.getUserByID)
	u.baseRouter.GET("/get", u.getUsersByID)
}

// createUser handles the POST /api/v1/users request
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body model.CreateUserInput true "User"
// @Success 201 {object} model.User
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /users [post]
func (u *UserRoutes) createUser(ctx *gin.Context) {
	u.userController.SetContext(ctx)
	var userInput model.CreateUserInput
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newUser := model.User{
		ID:          uuid.New().String(),
		Email:       userInput.Email,
		Username:    userInput.Username,
		DisplayName: userInput.DisplayName,
	}

	createdUser, err := u.userController.CreateUser(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, createdUser)
}

// getUserByID handles the GET /api/v1/users/:id request
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /users/{id} [get]
func (u *UserRoutes) getUserByID(ctx *gin.Context) {
	u.userController.SetContext(ctx)
	id := ctx.Param("id")
	user, err := u.userController.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, "User not found")
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// getUsersByID handles the GET /api/v1/users request
// @Summary Get a list of users by IDs
// @Description Get a list of users by IDs
// @Tags users
// @Accept  json
// @Produce  json
// @Param ids query []string true "User IDs"
// @Success 200 {array} model.User
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /users [get]
func (u *UserRoutes) getUsersByID(ctx *gin.Context) {
	u.userController.SetContext(ctx)
	ids := ctx.QueryArray("ids")
	users, err := u.userController.GetUsersByID(ids)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if len(users) == 0 {
		ctx.JSON(http.StatusNotFound, "Users not found")
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// updateUser handles the PATCH /api/v1/users request
// @Summary Update a user
// @Description Update a user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body model.User true "User"
// @Success 200 {object} model.User
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /users [patch]
func (u *UserRoutes) updateUser(ctx *gin.Context) {
	u.userController.SetContext(ctx)
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	updatedUser, err := u.userController.UpdateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, updatedUser)
}
