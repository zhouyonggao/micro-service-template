package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/mitchellh/mapstructure"
	"microServiceTemplate/api/demotemplateorder/v1"
	"microServiceTemplate/internal/biz"
	"microServiceTemplate/internal/pkg/metadatamanager"
)

type OrderService struct {
	v1.UnimplementedOrderServer
	orderBiz  *biz.OrderBiz
	logHelper *log.Helper
}

func NewOrderService(orderBiz *biz.OrderBiz, logHelper *log.Helper) *OrderService {
	return &OrderService{orderBiz: orderBiz, logHelper: logHelper}
}

func (o *OrderService) GetOrder(ctx context.Context, req *v1.ReqQueryOrder) (*v1.RespOrder, error) {
	_, err := metadatamanager.GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	//记录日志必须 用 withContext，否则追踪不上 userId
	//o.logHelper.WithContext(ctx).Error("模拟记录错误日志报错11")
	//o.logHelper.WithContext(ctx).Error("模拟记录错误日志报错22")
	//o.logHelper.WithContext(ctx).Error("模拟记录错误日志报错33")
	//o.logHelper.WithContext(ctx).Error("模拟记录错误日志报错44")
	//o.logHelper.WithContext(ctx).Error("模拟记录错误日志报错55")

	//fmt.Printf("开始 sleep")
	//time.Sleep(10 * time.Second)
	//fmt.Printf("结束 sleep")

	orderEnt, err := o.orderBiz.GetOrder(ctx, int(req.GetId()))
	//fmt.Println("业务请求完成：", orderEnt, err, userId)

	if err != nil {
		return nil, err
	}
	//转换模型，详情：https://pkg.go.dev/github.com/mitchellh/mapstructure
	resp := &v1.RespOrder{}
	_ = mapstructure.Decode(orderEnt, resp)
	return resp, nil
}
