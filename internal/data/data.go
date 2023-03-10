package data

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"microServiceTemplate/internal/conf"
	"microServiceTemplate/internal/data/ent"
)

// Data .
type Data struct {
	// Db 数据库client
	Db *ent.Client
	// Rdb Redis的client
	Rdb *redis.Client
	//Mq 的 producer
	MqProducer rocketmq.Producer
}

// NewData .
func NewData(c *conf.Bootstrap, hLog *log.Helper) (*Data, func(), error) {
	//构建 database
	db, errDb := buildDb(c.Data, hLog)
	//构建 redis
	rdb := buildRdb(c.Data, hLog)
	//构建 mq producer
	mqProducer, errMq := buildMqProducer(c.Data, hLog)

	cleanup := func() {
		if db != nil {
			if err := db.Close(); err != nil {
				hLog.Error("关闭 Db 失败：", err)
			}
		}

		if rdb != nil {
			if err := rdb.Close(); err != nil {
				hLog.Error("关闭 redis 失败：", err)
			}
		}

		if mqProducer != nil {
			if err := mqProducer.Shutdown(); err != nil {
				hLog.Error("关闭 rocketMQ producer 失败：", err)
			}
		}
		fmt.Println("关闭 data 中的连接资源已完成")
	}
	var resErr error = nil
	if errDb != nil {
		resErr = errDb
	} else if errMq != nil {
		resErr = errMq
	}

	return &Data{Db: db, Rdb: rdb, MqProducer: mqProducer}, cleanup, resErr
}

// GetDb 获取DB client，此方法会判断 context 是否有事务，如果有会返回带事务的 client
func (d Data) GetDb(ctx context.Context) *ent.Client {
	tx, ok := ctx.Value(ContextTxKey{}).(*ent.Tx)
	if ok {
		return tx.Client()
	}
	return d.Db
}

// buildDb 构建db client
func buildDb(c *conf.Data, hLog *log.Helper) (*ent.Client, error) {
	//生成 Db
	drv, err := sql.Open(c.Database.GetDriver(), c.Database.GetSource())
	if err != nil {
		fmt.Println("连接db 失败: ", err)
		return nil, err
	}
	drv.DB().SetMaxIdleConns(int(c.Database.MaxIdle))
	drv.DB().SetMaxOpenConns(int(c.Database.MaxOpen))
	drv.DB().SetConnMaxLifetime(c.Database.MaxLifetime.AsDuration())
	client := ent.NewClient(ent.Driver(drv))

	//todo：开发模式打开 debug
	client = client.Debug()
	return client, nil
}

// buildRdb 构建redis client
func buildRdb(c *conf.Data, hLog *log.Helper) *redis.Client {
	if c.Redis == nil {
		return nil
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:            c.Redis.Addr,
		Password:        c.Redis.Password,
		ReadTimeout:     c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout:    c.Redis.WriteTimeout.AsDuration(),
		PoolSize:        int(c.Redis.PoolSize),
		MinIdleConns:    int(c.Redis.MinIdleConns),
		ConnMaxIdleTime: c.Redis.ConnMaxIdleTime.AsDuration(),
		DB:              0,
	})
	//此时并没有发起连接，在使用时才会
	return rdb
}

// buildMqProducer rocket producer
func buildMqProducer(c *conf.Data, hLog *log.Helper) (rocketmq.Producer, error) {
	if c.RocketMq == nil {
		return nil, nil
	}
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(c.RocketMq.NameServers),
		producer.WithRetry(int(c.RocketMq.Retry)),
		producer.WithGroupName(c.RocketMq.ProductGroup),
		producer.WithDefaultTopicQueueNums(4), //如果不存在则默认4个
	)
	if err != nil {
		fmt.Println("创建 rocketMQ producer 失败: ", err)
		return nil, err
	}
	err = p.Start()
	if err != nil {
		fmt.Println("rocketMQ producer start 失败: ", err)
		return nil, err
	}

	//此时并没有发起连接，在使用时才会
	return p, nil
}
