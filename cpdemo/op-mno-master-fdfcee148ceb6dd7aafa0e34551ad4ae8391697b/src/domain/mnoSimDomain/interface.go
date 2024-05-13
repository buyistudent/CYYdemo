package mnoSimDomain

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-proto.git/gen"
	"context"
)

type (
	ConnectOperatorService interface {
		GetPackage(context.Context, MnoCpPackageDTO) (*gen.MnoCpPackageResp, error)
		OrderPackage(context.Context, MnoCpOrderDTO) (*gen.MnoCpOrderResp, error)
		StopSim(context.Context, MnoCpStopDTO) (*gen.MnoCpStopResp, error)
		ResumeSim(context.Context, MnoCpResumeDTO) (*gen.MnoCpResumeResp, error)
		GetUsage(context.Context, MnoCpUsageDTO) (*gen.MnoCpUsageResp, error)
		GetSimInfo(context.Context, MnoCpStatusDTO) (*gen.MnoGetSimInfoResp, error)
		SentMessage(context.Context, MnoCpSentDTO) (*gen.MnoCpSentResp, error)
		GetProductOrderList(context.Context, ProductOrderListDTO) (*gen.MnoCpProductOrderListResp, error)
		SmsDetails(context.Context, GetSmsByIdDTO) (*gen.MnoCpSmsDetailsResp, error)
		MnoCpSimTestOrder(context.Context, PlaceOrderDTO) (*gen.MnoCpUpdatSimTestResp, error)
		MnoCpSimSendMessage(context.Context, SimMessage) (*gen.MnoCpSendMessageResp, error)
		UpdateSimEffectiveTime(context.Context, SimEffectiveTimeDTO) error
	}

	SimRepo interface {
		SaveBillESim(context.Context, *ESimBill) (*ESimBill, error)
		//更新esim查询变更流水
		UpdateESimBill(context.Context, *ESimBill) error
	}
)
