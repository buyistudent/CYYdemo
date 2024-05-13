//go:build wireinject
// +build wireinject

package inject

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-pkg.git/mysql8"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/app/config"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/infras_adapter/grpc/client"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/infras_adapter/grpc/server"
	repo "codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/infras_adapter/repo"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/usecases/sim"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		sim.SimUseCaseSet,
		server.MnoSimGRPCServerSet,
		client.MnoConnectSet,
		repo.SimRepositorySet,
	))
}

func dbEngineFunc(cfg *config.Config) (mysql8.DBEngine, func(), error) {
	connStr := mysql8.DBConnString(cfg.DsnUrl)
	maxOpen := mysql8.DBPoolMax(cfg.PoolMax)
	maxIdle := mysql8.DBIdleConnMax(cfg.IdleConnMax)
	time := mysql8.DBMaxIdleTime(cfg.MaxIdleTime)
	db, err := mysql8.NewMysql8DB(connStr, maxOpen, maxIdle, time)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}
