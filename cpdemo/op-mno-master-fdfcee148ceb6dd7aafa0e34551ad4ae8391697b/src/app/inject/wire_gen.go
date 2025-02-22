// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package inject

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-pkg.git/mysql8"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/app/config"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/infras_adapter/grpc/client"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/infras_adapter/grpc/server"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/infras_adapter/repo"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/usecases/CPSim"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/usecases/sim"
	"google.golang.org/grpc"

)

// Injectors from wire.go:
func InitApp(cfg *config.Config, grpcServer *grpc.Server) (*App, func(), error) {
	connectOperatorService, err := client.NewConnectOperatorService(cfg)
	if err != nil {
		return nil, nil, err
	}
	dbEngine, cleanup, err := dbEngineFunc(cfg)
	if err != nil {
		return nil, nil, err
	}
	simRepo := repo.NewSimRepo(dbEngine)
	mnoSimCase := sim.NewSimService(cfg, connectOperatorService, simRepo)
	mnoCpServiceServer := server.NewMnoSimGRPCServer(grpcServer, cfg, mnoSimCase)
	app := New(cfg, mnoCpServiceServer, dbEngine)

	CPSimcase:=CPSim.NewSimService(cfg)
	CPService:=server.NewCPSimServer(grpcServer,cfg,CPSimcase)
	if CPService==nil {

	}

	return app, func() {
		cleanup()
	}, nil
}

// wire.go:

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
