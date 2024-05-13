package shared

import "fmt"

type OrderStatus int8

type OrderResultStatus int8

const (
	Product_Order_UnFinish = 1
	Product_Order_Finish   = 2
	Product_Order_Close    = 3
)

const (
	ResultStatus = 0
)

func (o OrderStatus) String() string {
	return fmt.Sprintf("%d", int(o))
}

// 订单明细状态
type OrderDetailStatus int8

const (
	//未生效
	UnEffect = 0
	//已生效
	IsEffect = 1
	//已完成
	IsFinish = 2
	//已关闭
	IsClose = 3
	//全部
	All = 4
)

func (e OrderDetailStatus) String() string {
	return fmt.Sprintf("%d", int(e))
}

// 订单子状态
type OrderDetailChildStatus int8

const (
	//无
	Stateless = 0
	//待服务
	WaitServer = 1
	//开始服务
	StartServer = 2
)

func (c OrderDetailChildStatus) String() string {
	return fmt.Sprintf("%d", int(c))
}

// 删除状态
type IsDelete int8

const (
	//未删除
	No = 0
	//已删除
	Yes = 1
)

func (i IsDelete) String() string {
	return fmt.Sprintf("%d", int(i))
}

// sim卡状态
type SimStatus int8

const (
	//未生效
	SimEffect = 1
	//可测试
	CanTest = 2
	//库存
	Storage = 3
	//激活
	Action = 4
	//停卡
	Stop = 5
)

func (s SimStatus) String() string {
	return fmt.Sprintf("%d", int(s))
}

// apn类型
type ApnType int8

const (
	//APN1
	One = 1
	//APN2
	Two = 2
)

func (a ApnType) String() string {
	return fmt.Sprintf("%d", int(a))
}

// apn状态
type ApnStatus int8

const (
	//关闭
	Close = 1
	//打开
	Open = 2
)

func (apnStatus ApnStatus) String() string {
	return fmt.Sprintf("%d", int(apnStatus))
}

// 资费计划使用类型
type UseType int8

const (
	//立即生效
	InTime = 1
	//延时生效
	Delay = 2
)

func (useType UseType) String() string {
	return fmt.Sprintf("%d", int(useType))
}

const (
	Customer_Detail_Status_UnSuccess = 1
	Customer_Detail_Status_Success   = 2
)

const (
	Internal_Status = 1
	External_Status = 2
)

const (
	Apn_Up_Close = 1
	Apn_Up_Open  = 2
	Apn_Up       = 0
)

const (
	Bill_Status_UnComplete = 1
	Bill_Status_Complete   = 2
	Bill_Status_Unusual    = 3
)
const (
	Change_Type_Customer = 1 //1.客户开卡
	Change_Type_Simba    = 2 //2.客户停卡
	Change_Type_Operator = 3 //3.辛巴变更资费计划

)

const (
	Change_Event_Samemonth_Effect = 1 //1.当月生效
	Change_Event_Samemonth_Stop   = 2 //2.停卡状态
	Change_Event_Time_Change      = 3 // 3.开卡状态
	Change_Event_Time_Change_Stop = 4 // 4关停
)

type AbleStatus int8

const (
	//启用
	Enable = 1
	//禁用
	Disable = 2
)

func (a AbleStatus) string() string {
	return fmt.Sprintf("%d", int(a))
}

type SimPushStatus int8

const (
	//未推送
	UnPush = 1
	//推送成功
	PushSuccess = 2
	//推送失败
	PushFailed = 3
)

func (s SimPushStatus) String() string {
	return fmt.Sprintf("%d", int(s))
}

type ProductType int8

const (
	//常规商品
	ProductType_Normal = 1
	//测试商品
	ProductType_Test = 2
)

func (p ProductType) String() string {
	return fmt.Sprintf("%d", int(p))
}

type OrderType int8

const (
	//常规商品订单
	OrderType_Normal = 1
	//测试商品订单
	OrderType_Test = 2
)

func (o OrderType) String() string {
	return fmt.Sprintf("%d", int(o))
}

func GetOrderType(productType ProductType) OrderType {
	var orderType OrderType
	switch productType {
	case ProductType_Normal:
		orderType = OrderType_Normal
	case ProductType_Test:
		orderType = OrderType_Test
	}
	return orderType
}

type PlanType int8

const (
	//初始化
	PlanType_Default = 0
	//常规APN资费计划
	PlanType_Normal = 1
	//测试APN资费计划
	PlanType_Test = 2
)

func (o PlanType) String() string {
	return fmt.Sprintf("%d", int(o))
}

const (
	UnStatus = 0
	Status   = 1
)

type IsAbnormal int8

const (
	IsAbnormal_NO  = 1
	IsAbnormal_YES = 2
)

func (o IsAbnormal) String() string {
	return fmt.Sprintf("%d", int(o))
}

// orderStatus  0：全部订单 1：待生效 2：已生效-待服务 3：已生效-开始服务 4：已完成 5：已关闭
func ApnOrderStatusConvert(orderStatus int32) (OrderDetailStatus, OrderDetailChildStatus) {
	var orderDetailStatus OrderDetailStatus
	var orderDetailChildStatus OrderDetailChildStatus
	switch orderStatus {
	case 0:
		orderDetailStatus = All
		orderDetailChildStatus = Stateless
	case 1:
		orderDetailStatus = UnEffect
		orderDetailChildStatus = Stateless
	case 2:
		orderDetailStatus = IsEffect
		orderDetailChildStatus = WaitServer
	case 3:
		orderDetailStatus = IsEffect
		orderDetailChildStatus = StartServer
	case 4:
		orderDetailStatus = IsFinish
		orderDetailChildStatus = Stateless
	case 5:
		orderDetailStatus = IsClose
		orderDetailChildStatus = Stateless
	}
	return orderDetailStatus, orderDetailChildStatus
}

type ImportStatus int8

const (
	//校验通过
	Import_CheckSuccess = 1
	//导入中
	Import_Ing = 2
	//导入完成
	Import_Success = 3
	//导入失败
	Import_Fail = 4

	//定时任务导入失败
	Import_Time_Fail = 5
)

func (a ImportStatus) string() string {
	return fmt.Sprintf("%d", int(a))
}

type ErrType int8

const (
	Err_System      = 1
	Err_Interaction = 2
	Err_Normal      = 3
)

func (a ErrType) string() string {
	return fmt.Sprintf("%d", int(a))
}

type ImportRecordStatus int8

const (
	//导入中
	Import_Record_Ing = 2
	//导入完成
	Import_Record_Success = 3
)

type AccountProductStatus int8

const (
	//结账日订单使用状态
	AccountProductStatus_Use = 1
	//结账日订单完成状态
	AccountProductStatus_Finish = 2
)

func (a AccountProductStatus) string() string {
	return fmt.Sprintf("%d", int(a))
}
