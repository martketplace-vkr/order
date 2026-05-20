package config

import (
	cart "github.com/martketplace-vkr/cart/pkg/api/grpc/v1"
	"github.com/martketplace-vkr/order/internal/app/cmp/inbox"
	"github.com/martketplace-vkr/order/internal/app/cmp/outbox"

	"github.com/martketplace-vkr/pkg/build/components/pgxsqlxcomponent"
	"github.com/martketplace-vkr/pkg/kafkaconnector"
	"github.com/martketplace-vkr/pkg/server/grpc"
)

type Config struct {
	Grpc     grpc.Config                 `validate:"required"`
	Postgres pgxsqlxcomponent.Config     `validate:"required"`
	Cart     cart.Config                 `validate:"required"`
	Inbox    inbox.Config                `validate:"required"`
	Outbox   outbox.Config               `validate:"required"`
	Kafka    kafkaconnector.ClientConfig `validate:"required"`
}
