package sim

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-proto.git/gen"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/constant"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/exception"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/shared"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/utils"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/app/config"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/domain/mnoSimDomain"
	"context"
	"fmt"
	"golang.org/x/exp/slog"
	"time"

	"github.com/google/wire"
)

var _ MnoSimCase = (*mnoSimService)(nil)

var SimUseCaseSet = wire.NewSet(NewSimService)

type mnoSimService struct {
	cfg    *config.Config
	client mnoSimDomain.ConnectOperatorService
	repo   mnoSimDomain.SimRepo
}

func NewSimService(cfg *config.Config, client mnoSimDomain.ConnectOperatorService, repo mnoSimDomain.SimRepo) MnoSimCase {
	return &mnoSimService{
		cfg:    cfg,
		client: client,
		repo:   repo,
	}
}

func (s *mnoSimService) GetPackage(ctx context.Context, dto mnoSimDomain.MnoCpPackageDTO) (*gen.MnoCpPackageResp, error) {
	slog.Info("grpc GetPackage dto:", dto)
	result, err := s.client.GetPackage(ctx, dto)
	if err != nil {
		result.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		result.Message = exception.OPERATION_FAILURE.Error()
		return result, err
	}
	return result, nil

}

func (s *mnoSimService) OrderPackage(ctx context.Context, dto mnoSimDomain.MnoCpOrderDTO) (*gen.MnoCpOrderResp, error) {
	slog.Info("grpc OrderPackage dto:", dto)
	result, err := s.client.OrderPackage(ctx, dto)
	if err != nil {
		result.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		result.Message = exception.OPERATION_FAILURE.Error()
		return result, err
	}
	return result, nil
}

func (s *mnoSimService) StopSim(ctx context.Context, dto mnoSimDomain.MnoCpStopDTO) (*gen.MnoCpStopResp, error) {
	slog.Info("grpc StopSim dto:", dto)
	result, err := s.client.StopSim(ctx, dto)
	if err != nil {
		result.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		result.Message = exception.OPERATION_FAILURE.Error()
		return result, err
	}

	return result, nil
}

func (s *mnoSimService) GetSimStatus(ctx context.Context, dto mnoSimDomain.MnoCpStatusDTO) (*gen.MnoCpStatusResp, error) {
	slog.Info("grpc GetSimStatus dto:", dto)
	resps := gen.MnoCpStatusResp{}
	simInfo, err := s.client.GetSimInfo(ctx, dto)
	if err != nil {
		simInfo.Code = exception.OPERATION_FAILURE.Code()
		simInfo.Message = exception.OPERATION_FAILURE.Error()
		return &resps, err
	}
	if simInfo.Code == 200 {
		if simInfo.Result.SimStatus == 1 {
			resps.Code = fmt.Sprint(exception.Sim_UnEffect.Code())
			resps.Message = exception.Sim_UnEffect.Error()
			return &resps, err
		}
	}
	var status string = ""
	if simInfo.Code == 200 {
		simStatus := simInfo.Result.SimStatus
		slog.Info("resps.GetResult()", simStatus)
		if simStatus == 1 {
			resps.Code = "103"
			resps.Message = "SIM卡识别码不存在"
			return &resps, err
		} else if simStatus == 2 {
			status = "TEST"
		} else if simStatus == 3 {
			status = "INVENTORY"
		} else if simStatus == 4 {
			status = "ACTIVATED"
		} else if simStatus == 5 {
			status = "DEACTIVATED"
		}
		model := gen.MnoStatusInfoResp{
			Status: status,
		}
		resps.Result = &model
	}
	resps.Code = fmt.Sprint(simInfo.Code)
	resps.Message = simInfo.Message
	return &resps, err
}

func (s *mnoSimService) GetUsage(ctx context.Context, dto mnoSimDomain.MnoCpUsageDTO) (*gen.MnoCpUsageResp, error) {
	slog.Info("grpc GetUsage dto:", dto)
	usage, err := s.client.GetUsage(ctx, dto)
	if err != nil {
		usage.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		usage.Message = exception.OPERATION_FAILURE.Error()
		return usage, err
	}
	return usage, err
}

func (s *mnoSimService) SentMesage(ctx context.Context, dto mnoSimDomain.MnoCpSentDTO) (*gen.MnoCpSentResp, error) {
	slog.Info("grpc SentMesage dto:", dto)
	result, err := s.client.SentMessage(ctx, dto)
	if err != nil {
		result.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		result.Message = exception.OPERATION_FAILURE.Error()
		return result, err
	}

	return result, nil
}

func (s *mnoSimService) ResumeSim(ctx context.Context, dto mnoSimDomain.MnoCpResumeDTO) (*gen.MnoCpResumeResp, error) {
	slog.Info("grpc ResumeSim dto:", dto)
	result, err := s.client.ResumeSim(ctx, dto)
	if err != nil {
		result.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		result.Message = exception.OPERATION_FAILURE.Error()
		return result, err
	}
	return result, nil
}

func (s *mnoSimService) GetProductOrderList(ctx context.Context, dto mnoSimDomain.ProductOrderListDTO) (*gen.MnoCpProductOrderListResp, error) {
	slog.Info("grpc GetProductOrderList dto:", dto)
	result, err := s.client.GetProductOrderList(ctx, dto)
	if err != nil {
		result.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		result.Message = exception.OPERATION_FAILURE.Error()
		return result, err
	}

	return result, nil
}

func (s *mnoSimService) SmsDetails(ctx context.Context, dto mnoSimDomain.GetSmsByIdDTO) (*gen.MnoCpSmsDetailsResp, error) {
	slog.Info("grpc SmsDetails dto:", dto)
	result, err := s.client.SmsDetails(ctx, dto)
	if err != nil {
		result.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		result.Message = exception.OPERATION_FAILURE.Error()
		return result, err
	}

	return result, nil
}

func (s *mnoSimService) MnoCpSimTestOrder(ctx context.Context, dto mnoSimDomain.MnoCpSimTestOrderDTO) (*gen.MnoCpUpdatSimTestResp, error) {
	slog.Info("grpc MnoCpSimTestOrder dto:", dto)
	resp := gen.MnoCpUpdatSimTestResp{}
	dtos := mnoSimDomain.MnoCpStatusDTO{
		Identifier: dto.Identifier,
		SimId:      dto.SimId,
	}
	simInfo, err := s.client.GetSimInfo(ctx, dtos)
	if simInfo.Code == exception.Success.Code() {
		if err != nil {
			resp.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
			resp.Message = exception.OPERATION_FAILURE.Error()
			return &resp, err
		}

		if simInfo.Result.SimStatus != shared.SimEffect {
			resp.Code = fmt.Sprint(exception.ICCID_STATUS_UN_EFFECTIVE.Code())
			resp.Message = exception.ICCID_STATUS_UN_EFFECTIVE.Error()
			resp.Result = ""
			return &resp, nil
		}

		if simInfo.Result.TestProductid == 0 {
			resp.Code = fmt.Sprint(exception.ICCID_STATUS_TEST_PAGE.Code())
			resp.Message = exception.ICCID_STATUS_TEST_PAGE.Error()
			resp.Result = ""
			return &resp, nil
		}

		dtoz := mnoSimDomain.PlaceOrderDTO{
			ProductId:     simInfo.Result.TestProductid,
			Quantity:      1,
			Iccid:         simInfo.Result.Iccid,
			Identifier:    dto.Identifier,
			FromOrderID:   fmt.Sprint(utils.NextUUId()),
			FromOrderName: "mno平台",
			ClientId:      simInfo.Result.ClientId,
		}
		sim, err := s.client.MnoCpSimTestOrder(ctx, dtoz)
		if err != nil {
			sim.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
			sim.Message = exception.OPERATION_FAILURE.Error()
			return sim, err
		}
		slog.Info("grpc MnoCpSimTestOrder  下测试单:", sim.Code)
		resp.Code = sim.Code
		resp.Message = sim.Message
	} else {
		resp.Code = fmt.Sprint(simInfo.Code)
		resp.Message = simInfo.Message
	}
	return &resp, err
}

func (s *mnoSimService) MnoCpESimTestOrder(ctx context.Context, req mnoSimDomain.MnoCpUpdateSimCanTestDTO) (*gen.MnoCpUpdatESimTestResp, error) {
	slog.Info("grpc MnoCpESimTestOrder dto:", req)
	resp := gen.MnoCpUpdatESimTestResp{}
	esim := mnoSimDomain.ESimBill{
		ID:         utils.NextUUId(),
		Iccid:      req.DeviceDetails.Iccid,
		Imsi:       req.DeviceDetails.Imsi,
		Msisdn:     req.DeviceDetails.Msisdn,
		Json:       "RequestId:" + req.RequestId + "Code:" + req.EventStatus.Code + "Description:" + req.EventStatus.Description + "ProfileSwitchState:" + req.ProfileSwitchState + "OwningCarrier:" + req.OwningCarrier,
		CreateTime: time.Now(),
	}
	_, err2 := s.repo.SaveBillESim(ctx, &esim)
	if err2 != nil {
		slog.Info("SaveBillESim", "保存ESIM流水报错")
	}

	dtos := mnoSimDomain.MnoCpStatusDTO{
		Identifier: constant.ICCID_IDENTIFIER,
		SimId:      req.DeviceDetails.Iccid,
	}
	simInfo, err := s.client.GetSimInfo(ctx, dtos)
	slog.Info("MnoCpESimTestOrder,simInfo,", simInfo)
	if err != nil {
		resp.Codes = constant.Code_Test_Failed
		simTest := mnoSimDomain.RequestrespTest(constant.Code_Test_err3, "Iccid "+constant.Code_Test_err_mes3)
		resp.Body = &simTest
		return &resp, nil
	}
	if simInfo.Code == exception.SIM_NOT_EXIST.Code() {
		resp.Codes = constant.Code_Test_Failed
		simTest := mnoSimDomain.RequestrespTest(constant.Code_Test_err3, "Invalid Parameter:IMSI "+constant.Code_Test_err_mes3)
		resp.Body = &simTest
		return &resp, nil
	}

	if simInfo.Result.Imsi != req.DeviceDetails.Imsi {
		resp.Codes = constant.Code_Test_Failed
		simTest := mnoSimDomain.RequestrespTest(constant.Code_Test_err3, "Invalid Parameter:IMSI "+constant.Code_Test_err_mes3)
		resp.Body = &simTest
		return &resp, nil
	}
	if simInfo.Result.Msisdn != req.DeviceDetails.Msisdn {
		resp.Codes = constant.Code_Test_Failed
		simTest := mnoSimDomain.RequestrespTest(constant.Code_Test_err3, "Invalid Parameter:Msisdn "+constant.Code_Test_err_mes3)
		resp.Body = &simTest
		return &resp, nil
	}

	if simInfo.Result.SimStatus != shared.SimEffect {
		// resp.Codes = constant.Code_Test_Failed
		// simTest := simDomain.RequestrespTest(constant.Code_Test_err1, "Invalid Parameter:ICCID  status 不是未生效 "+constant.Code_Test_err_mes1)
		// resp.Body = &simTest
		resp.Codes = constant.Code_success
		return &resp, nil
	}

	effecModel := mnoSimDomain.SimEffectiveTimeDTO{
		EffectiveTime: time.Now(),
		SimId:         req.DeviceDetails.Iccid,
	}

	err3 := s.client.UpdateSimEffectiveTime(ctx, effecModel)
	if err3 != nil {
		slog.Info("UpdateSimEffectiveTime", "保存UpdateSimEffectiveTime报错")
	}

	if simInfo.Result.TestProductid == 0 {
		slog.Info("SiM 测试包不存在", req.DeviceDetails.Iccid, exception.ICCID_STATUS_TEST_PAGE.Error())
		resp.Codes = fmt.Sprint(exception.Success.Code())
		return &resp, nil
	}

	dtoOrder := mnoSimDomain.PlaceOrderDTO{
		ProductId:     simInfo.Result.TestProductid,
		Quantity:      1,
		Iccid:         simInfo.Result.Iccid,
		Identifier:    constant.ICCID_IDENTIFIER,
		FromOrderID:   fmt.Sprint(utils.NextUUId()),
		FromOrderName: "mno平台",
		ClientId:      simInfo.Result.ClientId,
	}

	a, err1 := s.client.MnoCpSimTestOrder(ctx, dtoOrder)
	if err1 != nil {
		resp.Codes = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		return &resp, nil
	}
	esim.Status = fmt.Sprint(shared.Status)
	esim.EffectiveTime = time.Now()
	err6 := s.repo.UpdateESimBill(ctx, &esim)
	if err6 != nil {
		slog.Info("UpdateESimBill", "更新ESIM流水报错")
	}
	slog.Info("grpc MnoCpSimTestOrder dto:", a.Code)
	resp.Codes = fmt.Sprint(exception.Success.Code())
	return &resp, nil
}

func (s *mnoSimService) MnoCpSimSendMessage(ctx context.Context, dto mnoSimDomain.SimMessage) (*gen.MnoCpSendMessageResp, error) {
	slog.Info("grpc MnoCpSimSendMessage dto:", dto)
	esims, err := s.client.MnoCpSimSendMessage(ctx, dto)
	if err != nil {
		esims.Code = fmt.Sprint(exception.OPERATION_FAILURE.Code())
		return esims, err
	}
	return esims, err
}

func (s *mnoSimService) SaveBillESim(ctx context.Context, bill *mnoSimDomain.ESimBill) (*mnoSimDomain.ESimBill, error) {
	esim, err1 := s.repo.SaveBillESim(ctx, bill)
	if err1 != nil {
		return nil, err1
	}
	return esim, nil
}

func (s *mnoSimService) UpdateSimEffectiveTime(ctx context.Context, bill *mnoSimDomain.SimEffectiveTimeDTO) error {
	err := s.client.UpdateSimEffectiveTime(ctx, *bill)
	if err != nil {
		return err
	}
	return nil
}

func (s *mnoSimService) UpdateESimBill(ctx context.Context, bill *mnoSimDomain.ESimBill) error {
	err1 := s.repo.UpdateESimBill(ctx, bill)
	if err1 != nil {
		return err1
	}
	return nil
}
