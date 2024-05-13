package server

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-proto.git/gen"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/exception"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/app/config"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/domain/mnoSimDomain"
	mnosim "codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/usecases/sim"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/wire"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var _ gen.MnoCpServiceServer = (*mnoSimGRPCServer)(nil)

var MnoSimGRPCServerSet = wire.NewSet(NewMnoSimGRPCServer)

// 包含接口的所有服务
type mnoSimGRPCServer struct {
	gen.UnimplementedMnoCpServiceServer
	conn *grpc.ClientConn
	uc   mnosim.MnoSimCase
}

// 创建服务器实例
func NewMnoSimGRPCServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc mnosim.MnoSimCase,

) gen.MnoCpServiceServer {
	//连接grpc服务，获取客户端连接
	conn, err := grpc.Dial(cfg.Client.ConnectClientUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	svc := mnoSimGRPCServer{
		conn: conn,
		uc:   uc,
	}
	//创建一个服务器实例
	//将服务器实例注册到grpc服务器上
	gen.RegisterMnoCpServiceServer(grpcServer, &svc)
	return &svc
	//返回服务器实例的指针
}

func (s *mnoSimGRPCServer) MnoCp_Package(ctx context.Context, req *gen.MnoCpPackageReq) (*gen.MnoCpPackageResp, error) {

	marshal, _ := json.Marshal(*req)
	slog.Info("grpc  mno  Package req:", string(marshal))

	model := mnoSimDomain.MnoCpPackageDTO{
		PackageType: req.PackageType,
	}
	result, err := s.uc.GetPackage(ctx, model)
	if err != nil {
		result.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		result.Message = exception.OPERATION_FAILURE.Error()
		return result, nil
	}
	return result, nil
}

func (s *mnoSimGRPCServer) MnoCp_Order(ctx context.Context, req *gen.MnoCpOrderReq) (*gen.MnoCpOrderResp, error) {

	marshal, _ := json.Marshal(*req)
	slog.Info("grpc mno Order req:", string(marshal))
	resp := gen.MnoCpOrderResp{}

	if req.Identifier == "" || req.SimId == "" || req.PackageId == "" {
		resp.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		resp.Message = exception.PARAMETER_INVALID.Error()
		return &resp, nil
	}

	model := mnoSimDomain.MnoCpOrderDTO{
		Identifier: req.Identifier,
		SimId:      req.SimId,
		PackageId:  req.PackageId,
	}
	result, err := s.uc.OrderPackage(ctx, model)
	if err != nil {
		result.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		result.Message = exception.OPERATION_FAILURE.Error()
		return result, nil
	}
	return result, nil

}
func (s *mnoSimGRPCServer) MnoCp_Stop(ctx context.Context, req *gen.MnoCpStopReq) (*gen.MnoCpStopResp, error) {

	marshal, _ := json.Marshal(*req)
	slog.Info("grpc mno Stop req:", string(marshal))
	resp := gen.MnoCpStopResp{}

	if req.Identifier == "" || req.SimId == "" {
		resp.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		resp.Message = exception.PARAMETER_INVALID.Error()
		return &resp, nil
	}
	model := mnoSimDomain.MnoCpStopDTO{
		Identifier: req.Identifier,
		SimId:      req.SimId,
	}
	result, err := s.uc.StopSim(ctx, model)
	if err != nil {
		result.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		result.Message = exception.PARAMETER_INVALID.Error()
		return result, nil
	}
	return result, nil

}
func (s *mnoSimGRPCServer) MnoCp_Resume(ctx context.Context, req *gen.MnoCpResumeReq) (*gen.MnoCpResumeResp, error) {

	marshal, _ := json.Marshal(*req)
	slog.Info("grpc mno Resume req:", string(marshal))
	resp := gen.MnoCpResumeResp{}

	if req.Identifier == "" || req.SimId == "" {
		resp.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		resp.Message = exception.PARAMETER_INVALID.Error()
		resp.Result = ""
		return &resp, nil
	}

	model := mnoSimDomain.MnoCpResumeDTO{
		Identifier: req.Identifier,
		SimId:      req.SimId,
	}
	result, err := s.uc.ResumeSim(ctx, model)
	if err != nil {
		result.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		result.Message = exception.PARAMETER_INVALID.Error()
		return result, nil
	}
	return result, nil

}

func (s *mnoSimGRPCServer) MnoCp_Usage(ctx context.Context, req *gen.MnoCpUsageReq) (*gen.MnoCpUsageResp, error) {

	marshal, _ := json.Marshal(*req)
	slog.Info("grpc mno Usage req:", string(marshal))
	resp := gen.MnoCpUsageResp{}

	if req.Identifier == "" || req.SimId == "" {
		resp.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		resp.Message = exception.PARAMETER_INVALID.Error()
		return &resp, nil
	}

	model := mnoSimDomain.MnoCpUsageDTO{
		Identifier: req.Identifier,
		SimId:      req.SimId,
	}
	resps, err := s.uc.GetUsage(ctx, model)
	if err != nil {
		resps.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		resps.Message = exception.PARAMETER_INVALID.Error()
		return &resp, nil
	}
	return resps, nil

}

func (s *mnoSimGRPCServer) MnoCp_Status(ctx context.Context, req *gen.MnoCpStatusReq) (*gen.MnoCpStatusResp, error) {

	marshal, _ := json.Marshal(*req)
	slog.Info("grpc mno Status req:", string(marshal))
	resp := gen.MnoCpStatusResp{}

	if req.Identifier == "" || req.SimId == "" {
		resp.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		resp.Message = exception.PARAMETER_INVALID.Error()
		return &resp, nil
	}

	model := mnoSimDomain.MnoCpStatusDTO{
		Identifier: req.Identifier,
		SimId:      req.SimId,
	}
	resps, err := s.uc.GetSimStatus(ctx, model)
	slog.Info("resps", resps)
	if err != nil {
		resps.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		resps.Message = exception.PARAMETER_INVALID.Error()
		return &resp, nil
	}
	return resps, nil

}

func (s *mnoSimGRPCServer) MnoCp_OrderRecord(ctx context.Context, req *gen.MnoCpProductOrderListReq) (*gen.MnoCpProductOrderListResp, error) {
	marshal, _ := json.Marshal(*req)
	slog.Info("grpc mno ProductOrderList req:", string(marshal))
	resp := gen.MnoCpProductOrderListResp{}

	if req.Identifier == "" || req.SimId == "" {
		resp.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		resp.Message = exception.PARAMETER_INVALID.Error()
		return &resp, nil
	}
	model := mnoSimDomain.ProductOrderListDTO{
		Identifier: req.Identifier,
		SimId:      req.SimId,
	}
	result, err := s.uc.GetProductOrderList(ctx, model)
	if err != nil {
		result.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		result.Message = exception.PARAMETER_INVALID.Error()
		return result, nil
	}
	return result, nil

}
func (s *mnoSimGRPCServer) MnoCp_Sent(ctx context.Context, req *gen.MnoCpSentReq) (*gen.MnoCpSentResp, error) {
	marshal, _ := json.Marshal(*req)
	slog.Info("grpc mno Sent req:", string(marshal))
	resp := gen.MnoCpSentResp{}

	if req.Identifier == "" || req.SimId == "" || req.Message == "" {
		resp.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		resp.Message = exception.PARAMETER_INVALID.Error()
		return &resp, nil
	}
	model := mnoSimDomain.MnoCpSentDTO{
		Identifier:      req.Identifier,
		SimId:           req.SimId,
		Message:         req.Message,
		MessageEncoding: req.MessageEncoding,
		DataCoding:      req.DataCoding,
	}
	result, err := s.uc.SentMesage(ctx, model)
	if err != nil {
		result.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		result.Message = exception.PARAMETER_INVALID.Error()
		return result, nil
	}
	return result, nil

}

func (s *mnoSimGRPCServer) MnoCp_SmsDetails(ctx context.Context, req *gen.MnoCpSmsDetailsReq) (*gen.MnoCpSmsDetailsResp, error) {
	marshal, _ := json.Marshal(*req)
	slog.Info("grpc mno SmsDetails req:", string(marshal))
	resp := gen.MnoCpSmsDetailsResp{}

	if req.SmsId == "" {
		resp.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		resp.Message = exception.PARAMETER_INVALID.Error()
		return &resp, nil
	}
	model := mnoSimDomain.GetSmsByIdDTO{
		ID: req.SmsId,
	}
	result, err := s.uc.SmsDetails(ctx, model)
	if err != nil {
		result.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		result.Message = exception.PARAMETER_INVALID.Error()
		return result, nil
	}
	return result, nil

}

func (s *mnoSimGRPCServer) MnoCp_UpdatSimTest(ctx context.Context, req *gen.MnoCpUpdatSimTestReq) (*gen.MnoCpUpdatSimTestResp, error) {

	marshal, _ := json.Marshal(*req)
	slog.Info("grpc mno MnoCp_SimTest req:", string(marshal))
	resp := gen.MnoCpUpdatSimTestResp{}

	if req.Identifier == "" || req.SimId == "" {
		resp.Code = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		resp.Message = exception.PARAMETER_INVALID.Error()
		return &resp, nil
	}
	model := mnoSimDomain.MnoCpSimTestOrderDTO{
		Identifier: req.Identifier,
		SimId:      req.SimId,
	}
	respsOrder, err := s.uc.MnoCpSimTestOrder(ctx, model)
	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = exception.OPERATION_FAILURE.Error()
		return &resp, nil
	}
	resp.Code = respsOrder.Code
	resp.Message = respsOrder.Message
	return &resp, nil
}

func (s *mnoSimGRPCServer) MnoCp_UpdatESimTest(ctx context.Context, req *gen.MnoCpUpdatESimTestReq) (*gen.MnoCpUpdatESimTestResp, error) {

	marshal, _ := json.Marshal(*req)
	slog.Info("grpc mno MnoCp_UpdatEsimTest req:", string(marshal))
	resp := gen.MnoCpUpdatESimTestResp{}
	if len([]rune(req.DeviceDetails.Iccid)) == 19 {
		req.DeviceDetails.Iccid = req.DeviceDetails.Iccid + "F"
	}
	resp, _ = mnoSimDomain.UpdateSimCanTestCheck(req)
	if resp.Codes != "" {
		return &resp, nil
	}
	eventStatus := mnoSimDomain.EventStatus{
		Code:        req.EventStatus.Code,
		Description: req.EventStatus.Description,
	}
	deviceDetails := mnoSimDomain.DeviceDetails{
		Iccid:  req.DeviceDetails.Iccid,
		Imsi:   req.DeviceDetails.Imsi,
		Msisdn: req.DeviceDetails.Msisdn,
	}
	mnoCpUpdateSimCanTestDTO := mnoSimDomain.MnoCpUpdateSimCanTestDTO{
		RequestId:          req.RequestId,
		EventStatus:        &eventStatus,
		DeviceDetails:      &deviceDetails,
		ProfileSwitchState: req.ProfileSwitchState,
		OwningCarrier:      req.OwningCarrier,
	}

	respsOrder, err := s.uc.MnoCpESimTestOrder(ctx, mnoCpUpdateSimCanTestDTO)
	if err != nil {
		resp.Codes = fmt.Sprint(exception.PARAMETER_INVALID.Code())
		return &resp, nil
	}
	resp.Codes = respsOrder.Codes
	if respsOrder.Body != nil {
		resp.Body = respsOrder.Body
		marshal, _ = json.Marshal(resp)

	}
	slog.Info("grpc MnoCp_UpdatESimTest resp:", resp.Codes)
	return &resp, nil

}

func (o *mnoSimGRPCServer) MnoCp_SendMessage(ctx context.Context, request *gen.MnoCpSendMessageReq) (*gen.MnoCpSendMessageResp, error) {
	marshal, _ := json.Marshal(*request)
	slog.Info("grpc mno MnoCp_SimSendMessage req:", string(marshal))
	resp := gen.MnoCpSendMessageResp{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		slog.Info("获取metadata失败")
	}
	clientIds, ok := md["client_id"]
	simMessage, err := mnoSimDomain.NewSimMessage(clientIds, request.Identifier, lo.Map(request.Sims, func(item *gen.MnoCpSimId, _ int) *mnoSimDomain.Sim {
		return &mnoSimDomain.Sim{
			SimId: item.SimId,
		}
	}), request.Message, request.MessageEncoding, request.DataCoding)

	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = err.Error()
		return &resp, nil
	}
	resps, err := o.uc.MnoCpSimSendMessage(ctx, *simMessage)
	if err != nil {
		resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		resp.Message = err.Error()
		return resps, nil
	}
	slog.Info("SimSendMessage()->resp: ", resp.Message)
	return resps, nil
}

type myserver struct {
	gen.UnimplementedCPSimServiceServer
}

func (s *myserver) GetCustomerSimStatus(ctx context.Context, req *gen.GetCustomerSimInfoReq) (*gen.GetCustomerSimStatusResp, error) {
	fmt.Println("启用myserver")
	return &gen.GetCustomerSimStatusResp{Code: "001"}, nil
}
