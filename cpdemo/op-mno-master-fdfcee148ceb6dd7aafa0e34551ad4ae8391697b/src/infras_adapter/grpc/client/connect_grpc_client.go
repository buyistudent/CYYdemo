package client

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-proto.git/gen"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/exception"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/shared"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/utils"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/app/config"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/domain/mnoSimDomain"
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"strconv"
	"strings"
	"time"
)

var _ mnoSimDomain.ConnectOperatorService = (*SimCpClient)(nil)

type SimCpClient struct {
	conn *grpc.ClientConn
	cfg  *config.Config
	//uc   sim.MnoSimCase
}

var MnoConnectSet = wire.NewSet(NewConnectOperatorService)

func NewConnectOperatorService(cfg *config.Config) (mnoSimDomain.ConnectOperatorService, error) {
	conn, err := grpc.Dial(cfg.Client.ConnectClientUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &SimCpClient{
		cfg:  cfg,
		conn: conn,
	}, nil
}

func (v *SimCpClient) GetPackage(ctx context.Context, dto mnoSimDomain.MnoCpPackageDTO) (*gen.MnoCpPackageResp, error) {
	slog.Info("调用CONNECT服务，获取套餐包请求：", dto)
	c := gen.NewCpMnoServiceClient(v.conn)
	resp := gen.MnoCpPackageResp{}
	Client := v.GetClient(ctx)

	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)
	req := gen.MnoGetPackagesReq{
		PackageType: dto.PackageType,
	}
	result, err := c.Mno_GetPackages(ctx, &req)

	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, err
	}
	slog.Info("调用CONNECT服务，订购套餐包请求result：", result)
	if result.Code == 200 {

		var Packages []*gen.Package
		for _, page := range result.Result.Package {
			dtos := gen.Package{
				PackageId:   page.PackageId,
				PackageName: page.PackageName,
				Operator:    page.Operator,
				PackageType: page.PackageType,
				ServerTime:  int64(page.ServerTime),
			}
			Packages = append(Packages, &dtos)
		}

		resp = gen.MnoCpPackageResp{
			Code:    fmt.Sprint(result.Code),
			Message: result.Message,
			Result:  &gen.MnoPageInfoResp{Package: Packages},
		}
	}
	resp.Code = fmt.Sprint(result.Code)
	resp.Message = result.Message

	return &resp, nil
}

func (v *SimCpClient) OrderPackage(ctx context.Context, dto mnoSimDomain.MnoCpOrderDTO) (*gen.MnoCpOrderResp, error) {
	slog.Info("调用CONNECT服务，订购套餐包请求：", dto)
	c := gen.NewCpMnoServiceClient(v.conn)
	resp := gen.MnoCpOrderResp{}
	Client := v.GetClient(ctx)

	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)

	req := gen.MnoPlaceOrderReq{
		Identifier:    dto.Identifier,
		SimId:         dto.SimId,
		PackageId:     dto.PackageId,
		FromOrderID:   fmt.Sprint(utils.NextUUId()),
		FromOrderName: "mno平台",
	}
	result, err := c.Mno_PlaceOrder(ctx, &req)
	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, err
	}
	slog.Info("调用MNO服务，订购套餐包结果result：", result)
	resp.Code = fmt.Sprint(result.Code)
	resp.Message = result.Message
	if result.Result != nil {
		orderz := gen.MnoCpOrderz{
			OrderId: result.Result.OrderId,
		}
		resp.Result = &orderz
	}
	return &resp, nil
}

func (v *SimCpClient) StopSim(ctx context.Context, dto mnoSimDomain.MnoCpStopDTO) (*gen.MnoCpStopResp, error) {
	slog.Info("调用CONNECT服务，停卡服务请求：", dto)
	c := gen.NewCpMnoServiceClient(v.conn)
	resp := gen.MnoCpStopResp{}

	Client := v.GetClient(ctx)

	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)
	req := gen.MnoStopSimReq{
		Identifier: dto.Identifier,
		SimId:      dto.SimId,
	}
	result, err := c.Mno_StopSim(ctx, &req)
	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, err
	}
	slog.Info("调用MNO服务，停卡服务结果result：", result)
	resp.Code = fmt.Sprint(result.Code)
	resp.Message = result.Message
	resp.Result = result.Result
	return &resp, nil
}

func (v *SimCpClient) ResumeSim(ctx context.Context, dto mnoSimDomain.MnoCpResumeDTO) (*gen.MnoCpResumeResp, error) {

	slog.Info("调用MNO服务，恢复SIM卡服务请求：", dto)
	c := gen.NewCpMnoServiceClient(v.conn)
	resp := gen.MnoCpResumeResp{}
	Client := v.GetClient(ctx)

	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)
	req := gen.MnoResumeSimReq{
		Identifier: dto.Identifier,
		SimId:      dto.SimId,
	}
	result, err := c.Mno_ResumeSim(ctx, &req)
	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, err
	}
	slog.Info("调用MNO服务，恢复SIM卡服务结果result：", result)
	resp.Code = fmt.Sprint(result.Code)
	resp.Message = result.Message
	resp.Result = result.Result

	return &resp, nil
}

func (v *SimCpClient) GetUsage(ctx context.Context, dto mnoSimDomain.MnoCpUsageDTO) (*gen.MnoCpUsageResp, error) {
	slog.Info("调用MNO服务，用量查询请求：", dto)
	c := gen.NewCpMnoServiceClient(v.conn)
	resp := gen.MnoCpUsageResp{}
	Client := v.GetClient(ctx)

	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)

	req := gen.MnoGetSimInfoReq{
		Identifier: dto.Identifier,
		SimId:      dto.SimId,
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)
	result, err := c.Mno_GetSimInfo(ctx, &req)
	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, err
	}
	slog.Info("调用MNO服务，用量查询返回结果result：", result)
	if result.Code == 200 {
		if result.Result.SimStatus == 1 {
			resp.Code = fmt.Sprint(exception.Sim_UnEffect.Code())
			resp.Message = exception.Sim_UnEffect.Error()
			return &resp, err
		}

		var dataLimit int32
		var dataUsage int32
		var dataLeft int32
		var expireDate string

		var entertDataLimit int32
		var entertDataUsage int32
		var entertDataLeft int32
		var entertExpireDate string

		for _, apn := range result.Result.Apns {
			if apn.ApnType == 1 {
				dataLimit = apn.ApnUp
				dataUsage = apn.SimbaApn.OperatorApn.ApnUsage
				dataLeft = apn.ApnUp - apn.SimbaApn.OperatorApn.ApnUsage
				if dataLeft < 0 {
					dataLeft = 0
				}
				if apn.ApnEndTime != "" {
					expireDate = v.fmtTimes(apn.ApnEndTime)
				} else {
					expireDate = "20221126"
				}
			} else if apn.ApnType == 2 {
				entertDataLimit = apn.ApnUp
				entertDataUsage = apn.SimbaApn.OperatorApn.ApnUsage
				entertDataLeft = apn.ApnUp - apn.SimbaApn.OperatorApn.ApnUsage
				if entertDataLeft < 0 {
					entertDataLeft = 0
				}
				if apn.ApnEndTime != "" {
					entertExpireDate = v.fmtTimes(apn.ApnEndTime)
				} else {
					entertExpireDate = "20221126"
				}
			}

		}

		basic := gen.Basic{
			DataLimit:  strconv.Itoa(int(dataLimit)),
			DataUsage:  strconv.Itoa(int(dataUsage)),
			DataLeft:   strconv.Itoa(int(dataLeft)),
			ExpireDate: expireDate,
		}
		entertainment := gen.Entertainment{
			DataLimit:  strconv.Itoa(int(entertDataLimit)),
			DataUsage:  strconv.Itoa(int(entertDataUsage)),
			DataLeft:   strconv.Itoa(int(entertDataLeft)),
			ExpireDate: entertExpireDate,
		}

		results := gen.MnoUsageInfoResp{
			Iccid:         result.Result.Iccid,
			Imsi:          result.Result.Imsi,
			Msisdn:        result.Result.Msisdn,
			Basic:         &basic,
			Entertainment: &entertainment,
		}

		resp.Result = &results
	}
	resp.Code = strconv.FormatInt(int64(result.Code), 10)
	resp.Message = result.Message

	return &resp, nil
}

func (v *SimCpClient) GetSimInfo(ctx context.Context, dto mnoSimDomain.MnoCpStatusDTO) (*gen.MnoGetSimInfoResp, error) {
	slog.Info("调用COONET服务，Sim状态查询,请求：", dto)
	c := gen.NewCpMnoServiceClient(v.conn)

	resp := gen.MnoGetSimInfoResp{}
	Client := v.GetClient(ctx)

	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)
	reqs := gen.MnoGetSimInfoReq{
		Identifier: dto.Identifier,
		SimId:      dto.SimId,
	}
	result, err := c.Mno_GetSimInfo(ctx, &reqs)
	if err != nil {
		resp.Code = exception.OPERATION_FAILURE.Code()
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, err
	}
	slog.Info("调用MNO服务，Sim状态查询返回结果result：", result)
	if result.Code == exception.Success.Code() {
		result.Code = result.Code
		result.Message = result.Message
	}
	return result, nil
}

func (v *SimCpClient) SentMessage(ctx context.Context, dto mnoSimDomain.MnoCpSentDTO) (*gen.MnoCpSentResp, error) {
	slog.Info("调用MNO服务，发送短信至SIM卡请求：", dto)
	u := gen.NewCpMnoServiceClient(v.conn)
	resp := gen.MnoCpSentResp{}
	Client := v.GetClient(ctx)
	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)
	reqs := gen.MnoGetSimInfoReq{
		Identifier: dto.Identifier,
		SimId:      dto.SimId,
	}
	resultSimInfo, err := u.Mno_GetSimInfo(ctx, &reqs)
	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, err
	}
	if resultSimInfo.Code == exception.Success.Code() {
		if resultSimInfo.Result.SimStatus == 1 {
			resp.Code = fmt.Sprint(exception.Sim_UnEffect.Code())
			resp.Message = exception.Sim_UnEffect.Error()
			return &resp, err
		}

		//req := gen.MnoSendMessageReq{
		//	Identifier: dto.Identifier,
		//	Sims: append([]*gen.SimId{}, &gen.SimId{
		//		SimId: dto.SimId,
		//	}),
		//}

		req := gen.MnoSendMessageReq{
			Identifier:      dto.Identifier,
			SimId:           dto.SimId,
			Message:         dto.Message,
			MessageEncoding: dto.MessageEncoding,
			DataCoding:      dto.DataCoding,
		}
		result, err := u.Mno_SendMessage(ctx, &req)
		if err != nil {
			resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
			resp.Message = exception.OPERATION_FAILURE.Error()
			return &resp, err
		}
		slog.Info("调用MNO服务，发送短信至SIM卡返回结果result：", result)
		resp.Code = fmt.Sprint(result.Code)
		resp.Message = result.Message
		if result.Code == fmt.Sprint(exception.Success.Code()) {
			info := gen.SentResult{
				SmsId: result.Result.SmsId,
			}
			resp.Result = &info
		}
	}
	resp.Code = fmt.Sprint(resultSimInfo.Code)
	resp.Message = resultSimInfo.Message
	return &resp, nil
}

func (v *SimCpClient) GetProductOrderList(ctx context.Context, dto mnoSimDomain.ProductOrderListDTO) (*gen.MnoCpProductOrderListResp, error) {
	slog.Info("调用COONET服务，GetProductOrderList查询,请求：", dto)
	c := gen.NewCpMnoServiceClient(v.conn)
	resp := gen.MnoCpProductOrderListResp{}
	Client := v.GetClient(ctx)

	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)
	reqs := gen.MnoGetProductOrderListReq{
		Identifier: dto.Identifier,
		SimId:      dto.SimId,
	}
	result, err := c.Mno_GetProductOrderList(ctx, &reqs)
	if err != nil {
		resp.Code = string(exception.OPERATION_FAILURE.Code())
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, err
	}
	if result.Result != nil {

		var orders []*gen.CpOrder
		for _, list := range result.Result.Order {
			var isStatus string // 0  待生效 1.生效中（待完成）2.已完成
			if list.OrderStatus == 2 {
				dateTime := shared.FindDateAddmonths(shared.Format(list.FinishTime))
				result := shared.CompareData(dateTime, time.Now()) // -1 ：t1小于t2  0:t1等于t2  1：t1大于t2
				if result == -1 {
					isStatus = "2"
				} else {
					isStatus = "1"
				}
			} else if list.OrderStatus == 3 {
				isStatus = "2"
			} else if list.OrderStatus == 1 && list.EffectTime != "" {
				isStatus = "1"
			} else if list.OrderStatus == 1 && list.EffectTime == "" {
				isStatus = "0"
			}

			var finishTime string
			if list.FinishTime != "" {
				finishTime = v.fmtTimes(list.FinishTime)
			} else {
				finishTime = ""
			}

			var effectTime string
			if list.EffectTime != "" {
				effectTime = v.fmtTimes(list.EffectTime)
			} else {
				effectTime = ""
			}
			order := gen.CpOrder{
				OrderId:     list.OrderId,
				PackageId:   list.Product.ProductId,
				PackageName: list.Product.ProductName,
				ServerTime:  list.Product.ServerTime,
				OrderState:  isStatus,
				OrderTime:   v.fmtTimes(list.CreateTime),
				FinishTime:  finishTime,
				EffectTime:  effectTime,
			}
			orders = append(orders, &order)
		}
		cpProductOrder := gen.CpProductOrder{
			Order: orders,
		}
		resp.Result = &cpProductOrder
	}
	slog.Info("调用MNO服务，GetProductOrderList返回结果result：", result)
	resp.Code = fmt.Sprint(result.Code)
	resp.Message = result.Message
	return &resp, nil
}

func (v *SimCpClient) fmtTimes(times string) string {
	str := times
	layout := "2006-01-02 15:04:05"
	location, _ := time.Parse(layout, str)
	format := location.Format("2006-01-02")
	zime := strings.Replace(format, "-", "", -1)
	return zime
}

func (v *SimCpClient) GetClient(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		slog.Info("获取metadata失败")
	}
	clientIds, ok := md["client_id"]
	if !ok {
		slog.Info("获取client_id失败")
	}
	var Client string
	if len(clientIds) != 0 && clientIds[0] != "" {
		Client = clientIds[0]
	}
	return Client
}

func (v *SimCpClient) SmsDetails(ctx context.Context, dto mnoSimDomain.GetSmsByIdDTO) (*gen.MnoCpSmsDetailsResp, error) {
	slog.Info("调用MNO服务，获取短信详情请求：", dto)
	c := gen.NewCpMnoServiceClient(v.conn)
	resp := gen.MnoCpSmsDetailsResp{}
	Client := v.GetClient(ctx)
	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)
	req := gen.MnoGetSimSmsReq{
		SmsId: dto.ID,
	}
	result, err := c.Mno_GetSimSms(ctx, &req)
	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, err
	}
	slog.Info("调用MNO服务，发送短信至SIM卡返回结果result：", result)

	resp.Code = fmt.Sprint(result.Code)
	resp.Message = result.Message
	if result.Code == exception.Success.Code() {
		smsDetail := gen.MnoSimSmsDetail{
			SmsId:        result.Result.SmsId,
			Iccid:        result.Result.Iccid,
			Imsi:         result.Result.Imsi,
			Msisdn:       result.Result.Msisdn,
			SmsType:      result.Result.SmsType,
			SmsMsg:       result.Result.SmsMsg,
			DataCoding:   result.Result.DataCoding,
			SmsStatus:    result.Result.SmsStatus,
			SentTime:     result.Result.SentTime,
			ReceivedTime: result.Result.ReceivedTime,
		}
		resp.Result = &smsDetail
	}

	return &resp, nil
}

func (v *SimCpClient) MnoCpSimTestOrder(ctx context.Context, req mnoSimDomain.PlaceOrderDTO) (*gen.MnoCpUpdatSimTestResp, error) {
	slog.Info("调用COONET服务，MnoCpSimTestOrder,请求：", req)
	c := gen.NewCPOrderServiceClient(v.conn)
	resp := gen.MnoCpUpdatSimTestResp{}
	Client := v.GetClient(ctx)
	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)

	request := gen.PlaceOrderReq{
		ProductId:     fmt.Sprint(req.ProductId),
		Quantity:      1,
		SimIds:        req.Iccid,
		Identifier:    req.Identifier,
		FromOrderID:   fmt.Sprint(utils.NextUUId()),
		FromOrderName: "mno平台",
		ClientId:      req.ClientId,
	}

	result, err := c.PlaceOrder(ctx, &request)
	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, err
	}
	resp.Code = fmt.Sprint(result.Code)
	resp.Message = result.Message
	return &resp, nil
}

func (v *SimCpClient) UpdateSimEffectiveTime(ctx context.Context, req mnoSimDomain.SimEffectiveTimeDTO) error {
	slog.Info("调用COONET服务，UpdateSimEffectiveTime,请求：", req)
	c := gen.NewCPSimServiceClient(v.conn)
	Client := v.GetClient(ctx)
	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)
	request := gen.UpdateSimEffectiveTimeReq{
		EffectiveTime: shared.FormatTimeToStr(req.EffectiveTime),
		SimId:         req.SimId,
	}
	_, err := c.UpdateSimEffectiveTime(ctx, &request)
	if err != nil {
		return err
	}
	return nil
}

func (v *SimCpClient) MnoCpSimSendMessage(ctx context.Context, dto mnoSimDomain.SimMessage) (*gen.MnoCpSendMessageResp, error) {

	slog.Info("调用COONET服务，MnoCpSimSendMessage,请求：", dto)
	c := gen.NewCPSimServiceClient(v.conn)
	resp := gen.MnoCpSendMessageResp{}
	Client := v.GetClient(ctx)
	ctx = metadata.AppendToOutgoingContext(ctx, "client_id", Client)
	request := &gen.SendMessageReq{
		Identifier: dto.Identifier,
		Sims: lo.Map(dto.Sims, func(item *mnoSimDomain.Sim, _ int) *gen.SimId {
			return &gen.SimId{
				SimId: item.SimId,
			}
		}),
		Message:         dto.Message,
		MessageEncoding: dto.MessageEncoding,
		DataCoding:      dto.DataCoding,
	}

	simResp, err := c.SendMessage(ctx, request)
	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, err
	}
	//if simResp.Code != strconv.Itoa(int(exception.Success.Code())) {
	//	resp.Code = simResp.Code
	//	resp.Message = simResp.Message
	//} else {
	result := &gen.MnoCpSendMessageResp{
		Code:    simResp.Code,
		Message: simResp.Message,
		Result: &gen.MnoCpSendMessageResultInfo{
			Success: lo.Map(simResp.Result.Success, func(item *gen.SendMessageResult, _ int) *gen.MnoCpSendMessageResult {
				return &gen.MnoCpSendMessageResult{
					OpCode: item.OpCode,
					SimId:  item.SimId,
					Reason: item.Reason,
				}
			}),
			Failure: lo.Map(simResp.Result.Failure, func(item *gen.SendMessageResult, _ int) *gen.MnoCpSendMessageResult {
				return &gen.MnoCpSendMessageResult{
					OpCode: item.OpCode,
					SimId:  item.SimId,
					Reason: item.Reason,
				}
			}),
			SerialNumber: simResp.GetResult().GetSerialNumber(),
		},
	}

	return result, nil

	return &resp, nil
}
