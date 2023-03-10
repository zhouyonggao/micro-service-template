package convert

import (
	"microServiceTemplate/internal/biz/do"
	"microServiceTemplate/internal/biz/valueobj"
	"microServiceTemplate/internal/data/ent"
)

var OrderConv = orderConv{}

type orderConv struct {
}

// ToEntity 转换成实体
func (orderConv) ToEntity(po *ent.Order) *do.OrderDo {
	if po == nil {
		return nil
	}
	return &do.OrderDo{
		Id:          po.ID,
		ProductName: po.ProductName,
		Status:      valueobj.OrderStatusObj(po.Status),
	}
}

// ToEntities 转换成实体
func (o orderConv) ToEntities(pos []*ent.Order) []*do.OrderDo {
	if pos == nil {
		return nil
	}
	entities := make([]*do.OrderDo, len(pos))
	for k, po := range pos {
		entities[k] = o.ToEntity(po)
	}
	return entities
}
