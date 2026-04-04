package app

import (
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	cartorderpb "github.com/martketplace-vkr/cart/pkg/api/grpc/v1/order"

	"github.com/martketplace-vkr/order/config"
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

	"github.com/martketplace-vkr/pkg/build"
	"github.com/martketplace-vkr/pkg/build/components/pgxsqlxcomponent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(ctx context.Context, cfg *config.Config) error {
	pg := pgxsqlxcomponent.New(cfg.Postgres)

	cartConn, err := dialGRPC(ctx, cfg.Cart)
	if err != nil {
		return err
	}
	defer cartConn.Close()

	clientRepo := clientRepository.New(pg.DB, trmsqlx.DefaultCtxGetter)
	adminRepo := adminRepository.New(pg.DB, trmsqlx.DefaultCtxGetter)
	vendorRepo := vendorRepository.New(pg.DB, trmsqlx.DefaultCtxGetter)
	cartClient := cartorderpb.NewCartOrderServiceClient(cartConn)
	clientServ := clientService.New(clientRepo, cartClient, cfg.Cart.Timeout.Duration)
	adminServ := adminService.New(adminRepo)
	vendorServ := vendorService.New(vendorRepo)

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
		grpcServer,
	}

	app, err := build.NewApp(cmps)
	if err != nil {
		return err
	}

	return build.Run(ctx, app)
}

func dialGRPC(ctx context.Context, cfg config.GRPCClient) (*grpc.ClientConn, error) {
	dialCtx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
	defer cancel()

	return grpc.DialContext(
		dialCtx,
		cfg.Host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
}
