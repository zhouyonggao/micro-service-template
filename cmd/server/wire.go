//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"microServiceTemplate/internal/biz"
	"microServiceTemplate/internal/conf"
	"microServiceTemplate/internal/data"
	"microServiceTemplate/internal/data/eventimpl"
	"microServiceTemplate/internal/data/repositoryimpl"
	"microServiceTemplate/internal/pkg"
	log2 "microServiceTemplate/internal/pkg/log"
	"microServiceTemplate/internal/server"
	"microServiceTemplate/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, log.Logger, *log2.AccessLogger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSetServer, data.ProviderSetData, biz.ProviderSetBiz, service.ProviderSetService,
		newApp, pkg.NewLogHelper, eventimpl.ProviderEvent, repositoryimpl.ProviderRepositoryImpl,
	))
}
