package admin

import adminpb "github.com/martketplace-vkr/order/pkg/api/grpc/v1/admin"

type Handler struct {
	service service
	adminpb.UnimplementedOrderAdminServiceServer
}

func New(service service) *Handler {
	return &Handler{
		service: service,
	}
}
