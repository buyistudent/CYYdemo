package exception

/*
异常枚举定义
规则：
	1xxxx: 订单域异常 如：10001：参数错误
	2xxxx: sim卡域异常
	3xxxx: 客户域异常
后续自行添加
*/

type ErrorCode int

const (
	PARAMETER_INVALID ErrorCode = 1 + iota
	SIM_NOT_EXIST
	OPERATION_FAILURE
	All_SUCCESS
	TOTAL_FAILURE
	PARTIAL_SUCCESS
	CUSTOMER_NOT_EXIST
	SIM_IS_EXIST
	IMSI_IS_EXIST
	MSISDN_IS_EXIST
	GET_OPERATOR_FAILURE
	ICCID_CANNOT_MODIFIED
	ACCOUNT_DAY_FAILURE
	APN_FLOW_PARAM_INVALID
	INVALID_SIM_DATA
	SIM_CANNOT_REPLACED

	Success
	Failed
	Sim_Place_Order_Failure
	Sim_UnEffect
	Sim_Not_Subscribe_Service
	Sim_Data_Is_Used_Up
	SMSID_NOT_EXIST
	SMS_NOT_EXIST
	ICCID_STATUS_UN_EFFECTIVE
	ICCID_STATUS_TEST_PAGE
)

type ErrorModel struct {
	Code    uint32
	Message string
}

var errorCodes = [...]ErrorModel{
	ErrorModel{Code: 100, Message: "request parameter error"},
	ErrorModel{Code: 103, Message: "SIM card identifier does not exist"},
	ErrorModel{Code: 500, Message: "operation failure"},
	ErrorModel{Code: 200, Message: "all success"},
	ErrorModel{Code: 102, Message: "total failure"},
	ErrorModel{Code: 103, Message: "partial success"},
	ErrorModel{Code: 104, Message: "customer information does not exist"},
	ErrorModel{Code: 105, Message: "SIM is exist"},
	ErrorModel{Code: 106, Message: "imei is exist"},
	ErrorModel{Code: 107, Message: "msisdn is exist"},
	ErrorModel{Code: 108, Message: "operator acquisition failure"},
	ErrorModel{Code: 109, Message: "iccid cannot modified"},
	ErrorModel{Code: 110, Message: "account date cannot be greater than 27 or negative"},
	ErrorModel{Code: 111, Message: "Illegal flow measurement unit"},
	ErrorModel{Code: 112, Message: "Invalid sim data"},
	ErrorModel{Code: 113, Message: "sim cannot be replaced"},

	{Code: 200, Message: "Operation succeeded"},
	{Code: 500, Message: "Business exception"},
	{Code: 101, Message: "SIM card status does not support this package"},
	{Code: 117, Message: "SIM card is invalid"},
	{Code: 114, Message: "sim card no subscribe service"},
	{Code: 115, Message: "sim card data is used up"},

	{Code: 116, Message: "Carrier failed to send SMS"},
	{Code: 118, Message: "The smsId does not exist"},
	{Code: 119, Message: "SIM status is not inactive"},
	{Code: 120, Message: "SIM No test package"},
}

func (c ErrorCode) Code() uint32 {
	return errorCodes[c-1].Code
}
func (c ErrorCode) Error() string {
	return errorCodes[c-1].Message
}
