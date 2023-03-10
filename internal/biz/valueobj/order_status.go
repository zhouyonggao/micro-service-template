package valueobj

type OrderStatusObj int

const (
	OrderStatusInit OrderStatusObj = iota
	OrderStatusPaid
	OrderStatusCancel
)

var orderStatusNameMap = map[OrderStatusObj]string{
	OrderStatusInit:   "无效",
	OrderStatusPaid:   "已支付",
	OrderStatusCancel: "已取消",
}

func (os OrderStatusObj) GetName() string {
	name, isExist := orderStatusNameMap[os]
	if isExist {
		return name
	}
	return "未知"
}

func (os OrderStatusObj) GetValue() int {
	return int(os)
}
