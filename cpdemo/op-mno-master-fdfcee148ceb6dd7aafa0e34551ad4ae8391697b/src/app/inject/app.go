package inject

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-pkg.git/mysql8"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-proto.git/gen"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/app/config"
)

type App struct {
	Cfg    *config.Config
	Server gen.MnoCpServiceServer
	DB     mysql8.DBEngine
}

func New(
	cfg *config.Config,
	server gen.MnoCpServiceServer,
	db mysql8.DBEngine,

) *App {
	return &App{
		Cfg:    cfg,
		Server: server,
		DB:     db,
	}
}
