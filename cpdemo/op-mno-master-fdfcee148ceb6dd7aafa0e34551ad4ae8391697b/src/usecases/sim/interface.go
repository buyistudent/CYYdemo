package sim

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-proto.git/gen"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/domain/mnoSimDomain"
	"context"
)

type (
	MnoSimCase interface {
		GetPackage(context.Context, mnoSimDomain.MnoCpPackageDTO) (*gen.MnoCpPackageResp, error)
		OrderPackage(context.Context, mnoSimDomain.MnoCpOrderDTO) (*gen.MnoCpOrderResp, error)
		StopSim(context.Context, mnoSimDomain.MnoCpStopDTO) (*gen.MnoCpStopResp, error)
		ResumeSim(context.Context, mnoSimDomain.MnoCpResumeDTO) (*gen.MnoCpResumeResp, error)
		GetUsage(context.Context, mnoSimDomain.MnoCpUsageDTO) (*gen.MnoCpUsageResp, error)
		GetSimStatus(context.Context, mnoSimDomain.MnoCpStatusDTO) (*gen.MnoCpStatusResp, error)
		SentMesage(context.Context, mnoSimDomain.MnoCpSentDTO) (*gen.MnoCpSentResp, error)
		GetProductOrderList(context.Context, mnoSimDomain.ProductOrderListDTO) (*gen.MnoCpProductOrderListResp, error)
		SmsDetails(context.Context, mnoSimDomain.GetSmsByIdDTO) (*gen.MnoCpSmsDetailsResp, error)
		MnoCpSimTestOrder(context.Context, mnoSimDomain.MnoCpSimTestOrderDTO) (*gen.MnoCpUpdatSimTestResp, error)
		MnoCpESimTestOrder(context.Context, mnoSimDomain.MnoCpUpdateSimCanTestDTO) (*gen.MnoCpUpdatESimTestResp, error)
		MnoCpSimSendMessage(context.Context, mnoSimDomain.SimMessage) (*gen.MnoCpSendMessageResp, error)
		//GetCPSimStatus(req CPSimDomain.CPSimUsageDTO) (*gen.GetCustomerSimUsageResp, error)

	}
)
