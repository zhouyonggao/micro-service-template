package repositoryimpl

import (
	"entgo.io/ent/dialect/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"microServiceTemplate/internal/data"
	"microServiceTemplate/internal/data/ent"
	"testing"
)

// GetMockData 获取 mock 了Db的 data
func GetMockData(t *testing.T) (*data.Data, sqlmock.Sqlmock, *miniredis.Miniredis) {
	db, mockDb := GetMockDb(t)
	rds, mockRds := GetMockRds(t)
	d := &data.Data{
		Db:  db,
		Rdb: rds,
	}
	return d, mockDb, mockRds
}

func GetMockDb(t *testing.T) (*ent.Client, sqlmock.Sqlmock) {
	//mock db
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new 失败, %s", err)
	}
	drv := sql.NewDriver("mysql", sql.Conn{ExecQuerier: db})
	client := ent.NewClient(ent.Driver(drv))
	return client, mockDb
}

// GetMockRds 获取 mock 的 rdb
func GetMockRds(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mockRds, err := miniredis.Run()
	if err != nil {
		t.Fatalf("mock redis 错误：%s", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: mockRds.Addr(), // mock redis server的地址
	})
	return rdb, mockRds
}
