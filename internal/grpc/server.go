package name_service_grpc

import (
	"context"
	"template-grpc-go-service/internal/domain/models"

	name_service_v1 "github.com/LiveisFpv/template-proto/gen/go/proto/template"
	"google.golang.org/grpc"
)

type serverAPI struct {
	name_service_v1.UnimplementedTestRPCServer
	name_service Name_service
}

// Methods needed for handlers on Service
type Name_service interface {
	Get_NamebyID(
		ctx context.Context,
		name_id int,
	) (name *models.Name, err error)
	Get_All_Name(
		ctx context.Context,
		pagination *models.Pagination,
		filter []*models.Filter,
		orderby []*models.Sort,
	) (
		countries []*models.Name,
		new_pagination *models.Pagination,
		err error)
	Add_Name(
		ctx context.Context,
		name_title string,
		name_capital string,
		name_area string,
	) (name *models.Name, err error)
	Update_NamebyID(
		ctx context.Context,
		name *models.Name,
	) (err error)
	Delete_NamebyID(
		ctx context.Context,
		name_id int,
	) (name *models.Name, err error)
}

// It how constructor but not constructor:В
func Register(gRPCServer *grpc.Server, name_service Name_service) {
	name_service_v1.RegisterTestRPCServer(gRPCServer, &serverAPI{name_service: name_service})
}
