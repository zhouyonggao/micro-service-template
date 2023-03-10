package service

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

type OrderConsumerService struct {
	logHelper *log.Helper
}

func NewOrderConsumerService(logHelper *log.Helper) *OrderConsumerService {
	return &OrderConsumerService{logHelper: logHelper}
}

func (o *OrderConsumerService) Test(body string) error {
	fmt.Println("------------Receive New messages : ", body)
	//time.Sleep(5 * time.Second)
	fmt.Println("------------Receive handle ok : ", body)
	return nil // 回调函数
}
