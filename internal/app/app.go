package app

import (
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	cart "github.com/martketplace-vkr/cart/pkg/api/grpc/v1"

	"github.com/martketplace-vkr/order/config"
	inboxComponent "github.com/martketplace-vkr/order/internal/app/cmp/inbox"
	outboxComponent "github.com/martketplace-vkr/order/internal/app/cmp/outbox"
	"github.com/martketplace-vkr/order/internal/app/cmp/server"
	adminRepository "github.com/martketplace-vkr/order/internal/repository/pg/admin"
	clientRepository "github.com/martketplace-vkr/order/internal/repository/pg/client"
	vendorRepository "github.com/martketplace-vkr/order/internal/repository/pg/vendor"
	adminService "github.com/martketplace-vkr/order/internal/service/admin"
	clientService "github.com/martketplace-vkr/order/internal/service/client"
	vendorService "github.com/martketplace-vkr/order/internal/service/vendor"
	adminTransport "github.com/martketplace-vkr/order/internal/transport/grpc/v1/admin"
	clientTransport "github.com/martketplace-vkr/order/internal/transport/grpc/v1/client"
	vendorTransport "github.com/martketplace-vkr/order/internal/transport/grpc/v1/vendor"
	"github.com/martketplace-vkr/order/pkg/eventmapper"

	"github.com/martketplace-vkr/pkg/build"
	"github.com/martketplace-vkr/pkg/build/components/pgxsqlxcomponent"
	"github.com/martketplace-vkr/pkg/kafkaconnector"
	outboxclient "github.com/martketplace-vkr/pkg/outbox"
)

func Run(ctx context.Context, cfg *config.Config) error {
	pg := pgxsqlxcomponent.New(cfg.Postgres)

	kafkaClient := kafkaconnector.NewClient(cfg.Kafka)
	kafkaProducer := kafkaClient.NewSyncProducer()

	outboxCl, err := outboxclient.NewDefaultWithOptions(
		cfg.Outbox.Outbox,
		outboxclient.WithSqlxDB(pg.DB),
		outboxclient.WithKafkaProducer(kafkaProducer),
	)
	if err != nil {
		return err
	}

	outboxCmp := outboxComponent.New(cfg.Outbox, outboxCl)

	cartClient := cart.New(cfg.Cart)
	clientRepo := clientRepository.New(pg.DB, trmsqlx.DefaultCtxGetter)
	adminRepo := adminRepository.New(pg.DB, trmsqlx.DefaultCtxGetter)
	vendorRepo := vendorRepository.New(pg.DB, trmsqlx.DefaultCtxGetter)
	clientServ := clientService.New(clientRepo, cartClient, outboxCmp)
	adminServ := adminService.New(adminRepo)
	vendorServ := vendorService.New(vendorRepo, outboxCmp)
	inboxCmp := inboxComponent.New(
		cfg.Inbox,
		pg,
		kafkaClient,
		eventmapper.GetEventMapper(adminServ),
	)

	clientHandler := clientTransport.New(clientServ)
	adminHandler := adminTransport.New(adminServ)
	vendorHandler := vendorTransport.New(vendorServ)

	grpcServer := server.New(
		cfg.Grpc,
		adminHandler,
		clientHandler,
		vendorHandler,
	)

	cmps := build.Components{
		pg,
		inboxCmp,
		outboxCmp,
		cartClient,
		grpcServer,
	}

	app, err := build.NewApp(cmps)
	if err != nil {
		return err
	}

	return build.Run(ctx, app)
}
