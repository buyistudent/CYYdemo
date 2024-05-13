package CPSimDomain

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-proto.git/gen"
	"context"
)

type (
	Myoperator interface {
		GetUsage(ctx context.Context, dto CPSimUsageDTO) (*gen.GetCustomerSimUsageResp, error)
	}
)
