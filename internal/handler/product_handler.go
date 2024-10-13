package handler

import (
	"net/http"
	"strconv"

	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/service"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	service service.ProductService
}

func RegisterProductHandler(e *echo.Echo, svc service.ProductService) {
	handler := &ProductHandler{
		service: svc,
	}
	e.GET("products", handler.GetAllProducts)
	e.POST("products", handler.CreateProduct)
	e.GET("products/:id", handler.GetProduct)
	e.PUT("products/:id", handler.UpdateProduct)
	e.DELETE("products/:id", handler.DeleteProduct)
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProduct Create product
// @Summary      Create Product
// @Description  Create Product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param product body model.Product true "Product details"
// @Success      201  {object}  model.Response
// @Failure      400  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product model.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	ctx := c.Request().Context()
	if err := h.service.Create(ctx, &product); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, model.Response{Message: "success", Data: product})
}

// GetProduct handles fetching an product by ID
// @Summary Get an product by ID
// @Description Retrieve an product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object}  model.Response
// @Failure 400 {object}  model.Response
// @Failure 404 {object}  model.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}

	ctx := c.Request().Context()
	product, err := h.service.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Product not found"})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: product})
}

// GetAllProducts handles fetching all product
// @Summary Get all product
// @Description Retrieve all product in the system
// @Tags products
// @Produce json
// @Success 200 {object}  model.Response
// @Failure 500 {object}  model.Response
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	ctx := c.Request().Context()
	products, err := h.service.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: products})
}

// UpdateProduct handles updating an existing product
// @Summary Update an existing product
// @Description Update product details
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body model.Product true "Product details"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	var product model.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}
	product.ID = id

	ctx := c.Request().Context()
	if err := h.service.Update(ctx, &product); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: product})
}

// DeleteProduct handles deleting an product by ID
// @Summary Delete an product by ID
// @Description Remove an product from the system by its ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 204
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
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
