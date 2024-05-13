package CPSim

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-proto.git/gen"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/domain/CPSimDomain"
	"context"
)

type (
	CPSimCase interface {
		GetCPSimStatus(ctx context.Context, dto CPSimDomain.CPSimUsageDTO) (*gen.GetCustomerSimStatusResp, error)
	}
)
