package server

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-proto.git/gen"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/app/config"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/domain/CPSimDomain"
	mnosim "codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/usecases/CPSim"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/exp/slog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

/*
实现获取sim卡状态，在mno里有，创建一个服务器，调用mno提供的mnoCP-status服务
*/
type CPSimServiceServer struct {
	gen.UnimplementedCPSimServiceServer
	uc mnosim.CPSimCase
}

func NewCPSimServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc mnosim.CPSimCase,

) gen.CPSimServiceServer {
	svc := CPSimServiceServer{
		uc: uc,
	}
	gen.RegisterCPSimServiceServer(grpcServer, &svc)
	return &svc
}

func (s *CPSimServiceServer) GetCustomerSimStatus(ctx context.Context, req *gen.GetCustomerSimInfoReq) (*gen.GetCustomerSimStatusResp, error) {
	marshal, _ := json.Marshal(req)
	slog.Info("get sim status error", string(marshal))
	if req.SimId == "" || req.Identifier == "" {
		return &gen.GetCustomerSimStatusResp{Code: "500", Message: "查询条件为空"}, nil
	}
	model := CPSimDomain.CPSimUsageDTO{SimId: req.SimId, Identifier: req.Identifier}
	resps, err := s.uc.GetCPSimStatus(ctx, model)
	if err != nil {
		slog.Info("get sim status error", string(marshal))
	}
	return resps, err
}
