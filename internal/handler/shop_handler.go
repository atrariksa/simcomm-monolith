package handler

import (
	"net/http"
	"strconv"

	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/service"

	"github.com/labstack/echo/v4"
)

type ShopHandler struct {
	service service.ShopService
}

func RegisterShopHandler(e *echo.Echo, svc service.ShopService) {
	handler := &ShopHandler{
		service: svc,
	}
	e.GET("shops", handler.GetAllShops)
	e.POST("shops", handler.CreateShop)
	e.GET("shops/:id", handler.GetShop)
	e.PUT("shops/:id", handler.UpdateShop)
	e.DELETE("shops/:id", handler.DeleteShop)
}

func NewShopHandler(service service.ShopService) *ShopHandler {
	return &ShopHandler{service: service}
}

// CreateShop Create shop
// @Summary      Create Shop
// @Description  Create Shop
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param shop body model.Shop true "Shop details"
// @Success      201  {object}  model.Response
// @Failure      400  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /shops [post]
func (h *ShopHandler) CreateShop(c echo.Context) error {
	var shop model.Shop
	if err := c.Bind(&shop); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	ctx := c.Request().Context()
	if err := h.service.Create(ctx, &shop); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, model.Response{Message: "success", Data: shop})
}

// GetShop handles fetching an shop by ID
// @Summary Get an shop by ID
// @Description Retrieve an shop by its ID
// @Tags shops
// @Produce json
// @Param id path int true "Shop ID"
// @Success 200 {object}  model.Response
// @Failure 400 {object}  model.Response
// @Failure 404 {object}  model.Response
// @Router /shops/{id} [get]
func (h *ShopHandler) GetShop(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}

	ctx := c.Request().Context()
	shop, err := h.service.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Shop not found"})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: shop})
}

// GetAllShops handles fetching all shop
// @Summary Get all shop
// @Description Retrieve all shop in the system
// @Tags shops
// @Produce json
// @Success 200 {object}  model.Response
// @Failure 500 {object}  model.Response
// @Router /shops [get]
func (h *ShopHandler) GetAllShops(c echo.Context) error {
	ctx := c.Request().Context()
	shops, err := h.service.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: shops})
}

// UpdateShop handles updating an existing shop
// @Summary Update an existing shop
// @Description Update shop details
// @Tags shops
// @Accept json
// @Produce json
// @Param id path int true "Shop ID"
// @Param shop body model.Shop true "Shop details"
// @Success 200 {object} model.Shop
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /shops/{id} [put]
func (h *ShopHandler) UpdateShop(c echo.Context) error {
	var shop model.Shop
	if err := c.Bind(&shop); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}
	shop.ID = id

	ctx := c.Request().Context()
	if err := h.service.Update(ctx, &shop); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: shop})
}

// DeleteShop handles deleting an shop by ID
// @Summary Delete an shop by ID
// @Description Remove an shop from the system by its ID
// @Tags shops
// @Param id path int true "Shop ID"
// @Success 204
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /shops/{id} [delete]
func (h *ShopHandler) DeleteShop(c echo.Context) error {
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
