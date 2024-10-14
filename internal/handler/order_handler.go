package handler

import (
	"net/http"
	"strconv"

	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/service"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	service service.OrderService
}

func RegisterOrderHandler(e *echo.Echo, svc service.OrderService) {
	handler := &OrderHandler{
		service: svc,
	}
	e.GET("orders", handler.GetAllOrders)
	e.POST("orders", handler.CreateOrder)
	e.GET("orders/:id", handler.GetOrder)
	e.PUT("orders/:id", handler.UpdateOrder)
	e.DELETE("orders/:id", handler.DeleteOrder)
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder Create order
// @Summary      Create Order
// @Description  Create Order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param order body model.Order true "Order details"
// @Success      201  {object}  model.Response
// @Failure      400  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var order model.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	ctx := c.Request().Context()
	if err := h.service.Create(ctx, &order); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, model.Response{Message: "success", Data: order})
}

// GetOrder handles fetching an order by ID
// @Summary Get an order by ID
// @Description Retrieve an order by its ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object}  model.Response
// @Failure 400 {object}  model.Response
// @Failure 404 {object}  model.Response
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}

	ctx := c.Request().Context()
	order, err := h.service.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Order not found"})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: order})
}

// GetAllOrders handles fetching all order
// @Summary Get all order
// @Description Retrieve all order in the system
// @Tags orders
// @Produce json
// @Success 200 {object}  model.Response
// @Failure 500 {object}  model.Response
// @Router /orders [get]
func (h *OrderHandler) GetAllOrders(c echo.Context) error {
	ctx := c.Request().Context()
	orders, err := h.service.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: orders})
}

// UpdateOrder handles updating an existing order
// @Summary Update an existing order
// @Description Update order details
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param order body model.Order true "Order details"
// @Success 200 {object} model.Order
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /orders/{id} [put]
func (h *OrderHandler) UpdateOrder(c echo.Context) error {
	var order model.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}
	order.ID = id

	ctx := c.Request().Context()
	if err := h.service.Update(ctx, &order); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: order})
}

// DeleteOrder handles deleting an order by ID
// @Summary Delete an order by ID
// @Description Remove an order from the system by its ID
// @Tags orders
// @Param id path int true "Order ID"
// @Success 204
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(c echo.Context) error {
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
