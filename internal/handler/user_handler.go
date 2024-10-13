package handler

import (
	"net/http"
	"strconv"

	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/service"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service service.UserService
}

func RegisterUserHandler(e *echo.Echo, svc service.UserService) {
	handler := &UserHandler{
		service: svc,
	}
	e.GET("users", handler.GetAllUsers)
	e.POST("users", handler.CreateUser)
	e.GET("users/:id", handler.GetUser)
	e.PUT("users/:id", handler.UpdateUser)
	e.DELETE("users/:id", handler.DeleteUser)

	e.POST("ecommerce/login", handler.Login)
	e.POST("ecommerce/signup", handler.SignUp)
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser    Create user
// @Summary      Create User
// @Description  Create User
// @Tags         users
// @Accept       json
// @Produce      json
// @Param user body model.User true "User details"
// @Success      201  {object}  model.Response
// @Failure      400  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	ctx := c.Request().Context()
	if err := h.service.Create(ctx, &user); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, model.Response{Message: "success", Data: user})
}

// GetUser handles fetching an user by ID
// @Summary Get an user by ID
// @Description Retrieve an user by its ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object}  model.Response
// @Failure 400 {object}  model.Response
// @Failure 404 {object}  model.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}

	ctx := c.Request().Context()
	user, err := h.service.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "User not found"})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: user})
}

// GetAllUsers handles fetching all user
// @Summary Get all user
// @Description Retrieve all user in the system
// @Tags users
// @Produce json
// @Success 200 {object}  model.Response
// @Failure 500 {object}  model.Response
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := h.service.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: users})
}

// UpdateUser handles updating an existing user
// @Summary Update an existing user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body model.User true "User details"
// @Success 200 {object} model.User
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}
	user.ID = id

	ctx := c.Request().Context()
	if err := h.service.Update(ctx, &user); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: user})
}

// DeleteUser handles deleting an user by ID
// @Summary Delete an user by ID
// @Description Remove an user from the system by its ID
// @Tags users
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}

	ctx := c.Request().Context()
	if err := h.service.Delete(ctx, id); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusNoContent, model.Response{Message: "success"})
}

// SignUp        Register Customer
// @Summary      Sign Up
// @Description  Sign Up
// @Tags         users
// @Accept       json
// @Produce      json
// @Param customer body model.SignUpRequest true "Customer details"
// @Success      201  {object}  model.Response
// @Failure      400  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /ecommerce/signup [post]
func (h *UserHandler) SignUp(c echo.Context) (err error) {
	var req model.SignUpRequest
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = req.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: err.Error()})
	}

	err = h.service.SignUp(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, model.Response{Message: "success"})
}

// Login         Customer Login
// @Summary      Login
// @Description  Login
// @Tags         users
// @Accept       json
// @Produce      json
// @Param loginRequest body model.LoginRequest true "Login request"
// @Success      201  {object}  model.Response
// @Failure      400  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /ecommerce/login [post]
func (h *UserHandler) Login(c echo.Context) (err error) {
	var req model.LoginRequest
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = req.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	data, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: data})
}
