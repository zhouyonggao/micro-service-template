package repositoryimpl

import (
	"github.com/google/wire"
	"microServiceTemplate/internal/data"
)

// ProviderRepositoryImpl is data providers.
var ProviderRepositoryImpl = wire.NewSet(NewOrderRepoImpl, data.NewTransaction)
