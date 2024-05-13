package mnoSimDomain

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-proto.git/gen"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/constant"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
	"strings"
	"time"
)

type MnoCpPackageDTO struct {
	PackageType string `json:"package_type"`
}

type MnoCpOrderDTO struct {
	Identifier string `json:"identifier"`
	SimId      string `json:"sim_id"`
	PackageId  string `json:"package_id"`
}

type MnoCpStopDTO struct {
	Identifier string `json:"identifier"`
	SimId      string `json:"sim_id"`
}

type MnoCpResumeDTO struct {
	Identifier string `json:"identifier"`
	SimId      string `json:"sim_id"`
}

type MnoCpUsageDTO struct {
	Identifier string `json:"identifier"`
	SimId      string `json:"sim_id"`
}

type MnoCpStatusDTO struct {
	Identifier string `json:"identifier"`
	SimId      string `json:"sim_id"`
}

type MnoCpSentDTO struct {
	Identifier      string `json:"identifier"`
	SimId           string `json:"sim_id"`
	Message         string `json:"message"`
	MessageEncoding string `json:"message_encoding"`
	DataCoding      string `json:"data_coding"`
}

type StatusDTO struct {
	Status string `json:"status"`
}

type SimId struct {
	SimId string `json:"sim_id"`
}
type Package struct {
	PackageId string `json:"package_id"`
	//套餐名称
	PackageName string `json:"package_name"`
	//运营商
	Operator string `json:"operator"`
	//套餐包类型 0:加油包 1：续费包
	PackageType string `json:"package_type"`
	//服务时长
	ServerTime int32 `json:"server_time"`
}

type MnoPageInfoResp struct {
	Package []*Package `json:"package"`
}

type MnoPageageDTO struct {
	Code    string           `json:"code"`
	Message string           `json:"messgae"`
	Result  *MnoPageInfoResp `json:"result"`
}

type ProductOrderListDTO struct {
	Identifier string `json:"identifier"`
	SimId      string `json:"sim_id"`
}

type GetSmsByIdDTO struct {
	ID string `json:"id"`
}

type MnoCpSimTestDTO struct {
	Identifier string `json:"identifier"`
	SimId      string `json:"sim_id"`
}

type MnoCpESimTestDTO struct {
	Identifier string `json:"identifier"`
	SimId      string `json:"sim_id"`
}

type PlaceOrderDTO struct {
	ProductId     int64 //链接平台商品id
	Quantity      int32 // 购买商品数量
	Iccid         string
	Identifier    string
	FromOrderID   string
	FromOrderName string
	ClientId      string
}

type ESimBill struct {
	ID int64

	Iccid         string
	Imsi          string
	Msisdn        string
	Json          string
	CreateTime    time.Time //订单创建时间
	Status        string
	EffectiveTime time.Time
}

type SimEffectiveTimeDTO struct {
	EffectiveTime time.Time `json:"effective_time"`
	SimId         string    `json:"sim_id"`
}

type MnoCpSimTestOrderDTO struct {
	Identifier string `json:"identifier"`
	SimId      string `json:"sim_id"`
}

type MnoCpUpdateSimCanTestDTO struct {
	RequestId          string         `json:"request_id"`
	EventStatus        *EventStatus   `json:"event_status"`
	DeviceDetails      *DeviceDetails `json:"device_details"`
	ProfileSwitchState string         `json:"profile_switch_state"`
	OwningCarrier      string         `json:"owning_carrier"`
}

type EventStatus struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type DeviceDetails struct {
	Iccid  string `json:"iccid"`
	Imsi   string `json:"imsi"`
	Msisdn string `json:"msisdn"`
}

type SimMessage struct {
	ClientId        string `json:"clientId"`        // SIM 卡标识符。Enums: ICCID IMSI MSISDN Label: 自定义标签
	Identifier      string `json:"identifier"`      // SIM 卡标识符。Enums: ICCID IMSI MSISDN Label: 自定义标签
	Sims            []*Sim `json:"sims"`            // SIM 集合
	Message         string `json:"message"`         // 短信发送内容。消息可以用十六进制编码以二进制方式发送
	MessageEncoding string `json:"messageEncoding"` // (可选)使用的短信编码类型
	DataCoding      string `json:"dataCoding"`      // (可选)使用的数据编码类型
}
type Sim struct {
	SimId string `json:"simId"` // SIM 识别号
}

func NewSimMessage(clientIds []string, identifier string, sims []*Sim, message string, messageEncode, dataCoding string) (*SimMessage, error) {
	simMessage := &SimMessage{}
	//if clientIds == nil || len(clientIds) <= 0 || strings.TrimSpace(clientIds[0]) == "" {
	//	return simMessage, errors.New("client_id不能为空")
	//}
	//simMessage.ClientId = clientIds[0]

	if strings.TrimSpace(identifier) == "" {
		return simMessage, errors.New("Business exception, identifier Cannot be empty")
	}
	simMessage.Identifier = identifier

	if sims == nil || len(sims) <= 0 || strings.TrimSpace(sims[0].SimId) == "" {
		return simMessage, errors.New("Business exception, sims Cannot be empty")
	}
	simMessage.Sims = sims

	if strings.TrimSpace(message) == "" {
		return simMessage, errors.New("Business exception, message Cannot be empty")
	}
	simMessage.Message = message

	simMessage.MessageEncoding = messageEncode
	simMessage.DataCoding = dataCoding
	return simMessage, nil
}

func RequestrespTest(pram1 string, pram2 string) gen.MnoCpSimTestErrs {
	var s []*gen.MnoCpSimErrs
	SimErrs := gen.MnoCpSimErrs{}
	SimErrs.Code = pram1
	SimErrs.Description = pram2
	s = append(s, &SimErrs)
	simTest := gen.MnoCpSimTestErrs{}
	simTest.Errors = s
	return simTest
}

func UpdateSimCanTestCheck(req *gen.MnoCpUpdatESimTestReq) (gen.MnoCpUpdatESimTestResp, error) {
	resp := gen.MnoCpUpdatESimTestResp{}
	if req.RequestId == "" {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err2, "RequestId "+constant.Code_Test_err_mes2)
		resp.Body = &simTest
		return resp, nil
	}

	if len([]rune(req.RequestId)) > 64 {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err1, "Invalid Parameter:RequestId "+constant.Code_Test_err_mes1)
		resp.Body = &simTest
		return resp, nil
	}
	if req.EventStatus.Code == "" {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err2, "Code "+constant.Code_Test_err_mes2)
		resp.Body = &simTest
		return resp, nil
	}

	if len([]rune(req.EventStatus.Code)) > 10 {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err1, "Invalid Parameter:Code "+constant.Code_Test_err_mes1)
		resp.Body = &simTest
		return resp, nil
	}

	if req.EventStatus.Description == "" {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err2, "Description "+constant.Code_Test_err_mes2)
		resp.Body = &simTest
		return resp, nil
	}

	if len([]rune(req.EventStatus.Description)) > 255 {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err1, "Invalid Parameter:Description "+constant.Code_Test_err_mes1)
		resp.Body = &simTest
		return resp, nil
	}

	if req.DeviceDetails.Iccid == "" {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err2, "ICCID "+constant.Code_Test_err_mes2)
		resp.Body = &simTest
		return resp, nil
	}
	//isIccid := shared.IsNumber(req.DeviceDetails.Iccid)

	if len([]rune(req.DeviceDetails.Iccid)) > 20 {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err1, "Invalid Parameter:ICCID "+constant.Code_Test_err_mes1)
		resp.Body = &simTest
		return resp, nil
	}
	if req.DeviceDetails.Imsi == "" {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err2, "IMSI "+constant.Code_Test_err_mes2)
		resp.Body = &simTest
		return resp, nil
	}
	//isImsi := shared.IsNumber(req.DeviceDetails.Imsi)
	if len([]rune(req.DeviceDetails.Imsi)) > 20 {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err1, "Invalid Parameter:IMSI "+constant.Code_Test_err_mes1)
		resp.Body = &simTest
		return resp, nil
	}
	if req.DeviceDetails.Msisdn == "" {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err2, "MSISDN "+constant.Code_Test_err_mes2)
		resp.Body = &simTest
		return resp, nil
	}
	//isMsisdn := shared.IsNumber(req.DeviceDetails.Msisdn)

	if len([]rune(req.DeviceDetails.Msisdn)) > 20 {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err1, "Invalid Parameter:MSISDN "+constant.Code_Test_err_mes1)
		resp.Body = &simTest
		return resp, nil
	}

	if req.ProfileSwitchState == "" {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err2, "ProfileSwitchState "+constant.Code_Test_err_mes2)
		resp.Body = &simTest
		return resp, nil
	}

	if req.OwningCarrier == "" {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err2, "OwningCarrier "+constant.Code_Test_err_mes2)
		resp.Body = &simTest
		return resp, nil
	}

	slog.Info("grpc customerSimInfo customerSimInfo.ProfileSwitchState:", req.ProfileSwitchState)
	if req.ProfileSwitchState != "Completed" {
		resp.Codes = constant.Code_Test_Failed
		simTest := RequestrespTest(constant.Code_Test_err1, "Invalid Parameter:ProfileSwitchState "+constant.Code_Test_err_mes1)

		resp.Body = &simTest
		return resp, nil
	}
	return resp, nil
}
