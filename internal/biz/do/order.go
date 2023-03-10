package do

import "microServiceTemplate/internal/biz/valueobj"

type OrderDo struct {
	Id          int
	ProductName string
	Status      valueobj.OrderStatusObj
}
