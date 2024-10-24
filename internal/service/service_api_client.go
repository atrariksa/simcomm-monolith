package service

type serviceAPIClient struct {
	UserService      UserService
	ProductService   ProductService
	OrderService     OrderService
	ShopService      ShopService
	WarehouseService WarehouseService
}

var SvcClients *serviceAPIClient

func GetNewServiceAPIClient() *serviceAPIClient {
	return &serviceAPIClient{}
}

func (sac *serviceAPIClient) AddUserService(s UserService) {
	sac.UserService = s
}

func (sac *serviceAPIClient) AddProductService(s ProductService) {
	sac.ProductService = s
}

func (sac *serviceAPIClient) AddOrderService(s OrderService) {
	sac.OrderService = s
}

func (sac *serviceAPIClient) AddShopService(s ShopService) {
	sac.ShopService = s
}

func (sac *serviceAPIClient) AddWarehouseService(s WarehouseService) {
	sac.WarehouseService = s
}
