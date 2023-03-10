package data

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

// ProviderSetData is data providers.
var ProviderSetData = wire.NewSet(NewData)
