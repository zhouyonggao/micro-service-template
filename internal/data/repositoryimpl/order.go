package repositoryimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"microServiceTemplate/internal/biz/do"
	"microServiceTemplate/internal/biz/repository"
	d "microServiceTemplate/internal/data"
	"microServiceTemplate/internal/data/convert"
	"microServiceTemplate/internal/data/ent/order"
	"strconv"
	"time"
)

type OrderRepoImpl struct {
	data *d.Data
}

func NewOrderRepoImpl(data *d.Data) repository.OrderRepo {
	return &OrderRepoImpl{
		data: data,
	}
}

// FindByID 根据id 查询数据，不存在返回 nil
func (o *OrderRepoImpl) FindByID(ctx context.Context, id int) *do.OrderDo {
	po := o.data.GetDb(ctx).Order.Query().Where(order.ID(id)).FirstX(ctx)
	if po == nil {
		return nil
	}
	fmt.Println(o.data.Rdb.Get(context.Background(), "api-secret").Result())
	// po 转 entity
	return convert.OrderConv.ToEntity(po)
}

// UpdateTimeNow 根据id更新当前时间
func (o *OrderRepoImpl) UpdateTimeNow(ctx context.Context, id int, now time.Time) error {
	_, err := o.data.GetDb(ctx).Order.Update().SetUpdateTime(now).Where(order.ID(id)).Save(ctx)
	return err
}

// FindByIDByRedis 根据id 查询数据，不存在返回 nil
func (o *OrderRepoImpl) FindByIDByRedis(ctx context.Context, id int) (*do.OrderDo, error) {
	key := "test_" + strconv.Itoa(id)
	res, err := o.data.Rdb.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}
	var do *do.OrderDo
	err = json.Unmarshal(res, &do)
	if err != nil {
		return nil, err
	}
	return do, nil
}
