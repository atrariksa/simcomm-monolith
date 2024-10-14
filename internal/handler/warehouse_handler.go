package handler

import (
	"net/http"
	"strconv"

	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/service"

	"github.com/labstack/echo/v4"
)

type WarehouseHandler struct {
	service service.WarehouseService
}

func RegisterWarehouseHandler(e *echo.Echo, svc service.WarehouseService) {
	handler := &WarehouseHandler{
		service: svc,
	}
	e.GET("warehouses", handler.GetAllWarehouses)
	e.POST("warehouses", handler.CreateWarehouse)
	e.GET("warehouses/:id", handler.GetWarehouse)
	e.PUT("warehouses/:id", handler.UpdateWarehouse)
	e.DELETE("warehouses/:id", handler.DeleteWarehouse)

	e.GET("warehouse-stored-products", handler.GetAllWarehouseStoredProducts)
	e.POST("warehouse-stored-products", handler.CreateWarehouseStoredProduct)
	e.GET("warehouse-stored-products/:id", handler.GetWarehouseStoredProduct)
	e.PUT("warehouse-stored-products/:id", handler.UpdateWarehouseStoredProduct)
	e.DELETE("warehouse-stored-products/:id", handler.DeleteWarehouseStoredProduct)
}

func NewWarehouseHandler(service service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

// CreateWarehouse Create warehouse
// @Summary      Create Warehouse
// @Description  Create Warehouse
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param warehouse body model.Warehouse true "Warehouse details"
// @Success      201  {object}  model.Response
// @Failure      400  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /warehouses [post]
func (h *WarehouseHandler) CreateWarehouse(c echo.Context) error {
	var warehouse model.Warehouse
	if err := c.Bind(&warehouse); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	ctx := c.Request().Context()
	if err := h.service.Create(ctx, &warehouse); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, model.Response{Message: "success", Data: warehouse})
}

// GetWarehouse handles fetching an warehouse by ID
// @Summary Get an warehouse by ID
// @Description Retrieve an warehouse by its ID
// @Tags warehouses
// @Produce json
// @Param id path int true "Warehouse ID"
// @Success 200 {object}  model.Response
// @Failure 400 {object}  model.Response
// @Failure 404 {object}  model.Response
// @Router /warehouses/{id} [get]
func (h *WarehouseHandler) GetWarehouse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}

	ctx := c.Request().Context()
	warehouse, err := h.service.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Warehouse not found"})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: warehouse})
}

// GetAllWarehouses handles fetching all warehouse
// @Summary Get all warehouse
// @Description Retrieve all warehouse in the system
// @Tags warehouses
// @Produce json
// @Success 200 {object}  model.Response
// @Failure 500 {object}  model.Response
// @Router /warehouses [get]
func (h *WarehouseHandler) GetAllWarehouses(c echo.Context) error {
	ctx := c.Request().Context()
	warehouses, err := h.service.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: warehouses})
}

// UpdateWarehouse handles updating an existing warehouse
// @Summary Update an existing warehouse
// @Description Update warehouse details
// @Tags warehouses
// @Accept json
// @Produce json
// @Param id path int true "Warehouse ID"
// @Param warehouse body model.Warehouse true "Warehouse details"
// @Success 200 {object} model.Warehouse
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /warehouses/{id} [put]
func (h *WarehouseHandler) UpdateWarehouse(c echo.Context) error {
	var warehouse model.Warehouse
	if err := c.Bind(&warehouse); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}
	warehouse.ID = id

	ctx := c.Request().Context()
	if err := h.service.Update(ctx, &warehouse); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: warehouse})
}

// DeleteWarehouse handles deleting an warehouse by ID
// @Summary Delete an warehouse by ID
// @Description Remove an warehouse from the system by its ID
// @Tags warehouses
// @Param id path int true "Warehouse ID"
// @Success 204
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /warehouses/{id} [delete]
func (h *WarehouseHandler) DeleteWarehouse(c echo.Context) error {
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

// CreateWarehouseStoredProduct Create warehousestoredproduct
// @Summary      Create WarehouseStoredProduct
// @Description  Create WarehouseStoredProduct
// @Tags         warehousestoredproducts
// @Accept       json
// @Produce      json
// @Param warehousestoredproduct body model.WarehouseStoredProduct true "WarehouseStoredProduct details"
// @Success      201  {object}  model.Response
// @Failure      400  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /warehouse-stored-products [post]
func (h *WarehouseHandler) CreateWarehouseStoredProduct(c echo.Context) error {
	var warehousestoredproduct model.WarehouseStoredProduct
	if err := c.Bind(&warehousestoredproduct); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	ctx := c.Request().Context()
	if err := h.service.WSPCreate(ctx, &warehousestoredproduct); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, model.Response{Message: "success", Data: warehousestoredproduct})
}

// GetWarehouseStoredProduct handles fetching an warehousestoredproduct by ID
// @Summary Get an warehousestoredproduct by ID
// @Description Retrieve an warehousestoredproduct by its ID
// @Tags warehousestoredproducts
// @Produce json
// @Param id path int true "WarehouseStoredProduct ID"
// @Success 200 {object}  model.Response
// @Failure 400 {object}  model.Response
// @Failure 404 {object}  model.Response
// @Router /warehouse-stored-products/{id} [get]
func (h *WarehouseHandler) GetWarehouseStoredProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}

	ctx := c.Request().Context()
	warehousestoredproduct, err := h.service.WSPGet(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "WarehouseStoredProduct not found"})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: warehousestoredproduct})
}

// GetAllWarehouseStoredProducts handles fetching all warehousestoredproduct
// @Summary Get all warehousestoredproduct
// @Description Retrieve all warehousestoredproduct in the system
// @Tags warehousestoredproducts
// @Produce json
// @Success 200 {object}  model.Response
// @Failure 500 {object}  model.Response
// @Router /warehouse-stored-products [get]
func (h *WarehouseHandler) GetAllWarehouseStoredProducts(c echo.Context) error {
	ctx := c.Request().Context()
	warehousestoredproducts, err := h.service.WSPGetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: warehousestoredproducts})
}

// UpdateWarehouseStoredProduct handles updating an existing warehousestoredproduct
// @Summary Update an existing warehousestoredproduct
// @Description Update warehousestoredproduct details
// @Tags warehousestoredproducts
// @Accept json
// @Produce json
// @Param id path int true "WarehouseStoredProduct ID"
// @Param warehousestoredproduct body model.WarehouseStoredProduct true "WarehouseStoredProduct details"
// @Success 200 {object} model.WarehouseStoredProduct
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /warehouse-stored-products/{id} [put]
func (h *WarehouseHandler) UpdateWarehouseStoredProduct(c echo.Context) error {
	var warehousestoredproduct model.WarehouseStoredProduct
	if err := c.Bind(&warehousestoredproduct); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}
	warehousestoredproduct.ID = id

	ctx := c.Request().Context()
	if err := h.service.WSPUpdate(ctx, &warehousestoredproduct); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Message: "success", Data: warehousestoredproduct})
}

// DeleteWarehouseStoredProduct handles deleting an warehousestoredproduct by ID
// @Summary Delete an warehousestoredproduct by ID
// @Description Remove an warehousestoredproduct from the system by its ID
// @Tags warehousestoredproducts
// @Param id path int true "WarehouseStoredProduct ID"
// @Success 204
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /warehouse-stored-products/{id} [delete]
func (h *WarehouseHandler) DeleteWarehouseStoredProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{Message: "Invalid ID"})
	}

	ctx := c.Request().Context()
	if err := h.service.WSPDelete(ctx, id); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusNoContent, model.Response{Message: "success"})
}
