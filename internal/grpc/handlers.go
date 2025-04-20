package name_service_grpc

import (
	"context"

	name_service_v1 "github.com/LiveisFpv/template-proto/gen/go/proto/template"
)

func (s *serverAPI) TestConnection(
	ctx context.Context, req *name_service_v1.TestConnectionRequest) (
	resp *name_service_v1.TestConnectionResponse, err error) {
	// name, err := s.name_service.Get_NamebyID(ctx, id)
	return &name_service_v1.TestConnectionResponse{
		TestMessage: req.TestMessage,
		TestValue:   req.TestValue,
	}, nil
}
