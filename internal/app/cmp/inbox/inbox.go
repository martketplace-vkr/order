package inbox

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/build/components/pgxsqlxcomponent"
	"github.com/martketplace-vkr/pkg/inbox"
	"github.com/martketplace-vkr/pkg/inbox/dto"
	"github.com/martketplace-vkr/pkg/kafkaconnector"
)

const cmpName = "inbox"

type cmp struct {
	cfg         Config
	pg          *pgxsqlxcomponent.PgxSqlxConnector
	kafkaClient kafkaconnector.Client
	inboxClient inbox.Inbox
	handlerMap  map[string]func(context.Context, dto.Event) error
}

func New(
	cfg Config,
	pg *pgxsqlxcomponent.PgxSqlxConnector,
	kafkaClient kafkaconnector.Client,
	handlerMap map[string]func(context.Context, dto.Event) error,
) *cmp {
	return &cmp{
		cfg:         cfg,
		pg:          pg,
		kafkaClient: kafkaClient,
		handlerMap:  handlerMap,
	}
}

func (c *cmp) Start(ctx context.Context) error {
	var err error
	c.inboxClient, err = inbox.NewDefaultWithOptions(
		c.cfg.Inbox,
		inbox.WithKafkaEventsProvider(c.kafkaClient.NewSaramaConsumerGroup),
		inbox.WithSqlxDB(c.pg.DB),
		inbox.WithHandlerMap(c.handlerMap),
	)
	if err != nil {
		return err
	}

	return c.inboxClient.Run(ctx)
}

func (c *cmp) Stop(ctx context.Context) error {
	if c.inboxClient == nil {
		return nil
	}
	return c.inboxClient.Shutdown(ctx)
}

func (c *cmp) GetName() string {
	return cmpName
}

func (c *cmp) GetShutdownDelay() time.Duration {
	return time.Second
}

func (c *cmp) GetStartTimeout() time.Duration {
	return 5 * time.Second
}

func (c *cmp) GetStopTimeout() time.Duration {
	return 5 * time.Second
}
