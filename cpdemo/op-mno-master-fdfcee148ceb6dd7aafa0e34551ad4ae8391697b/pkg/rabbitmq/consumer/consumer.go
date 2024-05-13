package consumer

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/rabbitmq"
	"context"

	"github.com/google/wire"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
)

const (
	_exchangeKind       = "direct"
	_exchangeDurable    = true
	_exchangeAutoDelete = false
	_exchangeInternal   = false
	_exchangeNoWait     = false

	_queueDurable    = true
	_queueAutoDelete = false
	_queueExclusive  = false
	_queueNoWait     = false

	_prefetchCount  = 5
	_prefetchSize   = 0
	_prefetchGlobal = false

	_consumeAutoAck   = false
	_consumeExclusive = false
	_consumeNoLocal   = false
	_consumeNoWait    = false

	_exchangeName   = "orders-exchange"
	_queueName      = "orders-queue"
	_bindingKey     = "orders-routing-key"
	_consumerTag    = "orders-consumer"
	_workerPoolSize = 24
)

type Consumer struct {
	exchangeName, queueName, bindingKey, consumerTag string
	workerPoolSize                                   int
	channelPool                                      *rabbitmq.ChannelPool
}

var EventConsumerSet = wire.NewSet(NewConsumer)

func NewConsumer(cp *rabbitmq.ChannelPool) (Consumer, error) {
	sub := Consumer{
		channelPool:    cp,
		exchangeName:   _exchangeName,
		queueName:      _queueName,
		bindingKey:     _bindingKey,
		consumerTag:    _consumerTag,
		workerPoolSize: _workerPoolSize,
	}

	return sub, nil
}

func (c *Consumer) Configure(opts ...Option) EventConsumer {
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// StartConsumer Start new rabbitmq consumer.
func (c *Consumer) StartConsumer(fn Worker) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.createChannel()
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	//defer ch.Close()

	deliveries, err := ch.Consume(
		c.queueName,
		c.consumerTag,
		_consumeAutoAck,
		_consumeExclusive,
		_consumeNoLocal,
		_consumeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Consume")
	}

	forever := make(chan bool)

	for i := 0; i < c.workerPoolSize; i++ {
		go fn(ctx, deliveries)
	}
	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	slog.Error("ch.NotifyClose", chanErr)
	if chanErr != nil {
		c.channelPool.CheckConnect()
		c.retryBindConsumer(ctx, fn, forever)
	}
	<-forever
	return amqp.ErrClosed
}
func (c *Consumer) retryBindConsumer(ctx context.Context, fn Worker, forever chan bool) {
	slog.Info("ch.NotifyClose--重启消费者")
	ch, err := c.createChannel()
	if err != nil {
		return
	}
	//defer ch.Close()

	deliveries, err := ch.Consume(
		c.queueName,
		c.consumerTag,
		_consumeAutoAck,
		_consumeExclusive,
		_consumeNoLocal,
		_consumeNoWait,
		nil,
	)
	if err != nil {
		return
	}
	for i := 0; i < c.workerPoolSize; i++ {
		go fn(ctx, deliveries)
	}
	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	slog.Error("ch.NotifyClose", chanErr)
	if chanErr != nil {
		c.channelPool.CheckConnect()
		c.retryBindConsumer(ctx, fn, forever)
	}
	<-forever
}

// CreateChannel Consume messages.
func (c *Consumer) createChannel() (*amqp.Channel, error) {
	ch, err := c.channelPool.GetChannel()
	if err != nil {
		return nil, errors.Wrap(err, "Error amqpConn.Channel")
	}

	slog.Info("declaring exchange", "exchange_name", c.exchangeName)
	err = ch.ExchangeDeclare(
		c.exchangeName,
		_exchangeKind,
		_exchangeDurable,
		_exchangeAutoDelete,
		_exchangeInternal,
		_exchangeNoWait,
		nil,
	)

	if err != nil {
		ch.Close()
		c.channelPool.CheckConnect()
		return nil, errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(
		c.queueName,
		_queueDurable,
		_queueAutoDelete,
		_queueExclusive,
		_queueNoWait,
		nil,
	)
	if err != nil {
		ch.Close()
		c.channelPool.CheckConnect()
		return nil, errors.Wrap(err, "Error ch.QueueDeclare")
	}

	slog.Info("declared queue, binding it to exchange", "queue", queue.Name, "messages_count", queue.Messages,
		"consumer_count", queue.Consumers, "exchange", c.exchangeName, "binding_key", c.bindingKey,
	)

	err = ch.QueueBind(
		queue.Name,
		c.bindingKey,
		c.exchangeName,
		_queueNoWait,
		nil,
	)
	if err != nil {
		ch.Close()
		c.channelPool.CheckConnect()
		return nil, errors.Wrap(err, "Error ch.QueueBind")
	}

	slog.Info("queue bound to exchange, starting to consume from queue", "consumer_tag", c.consumerTag)

	err = ch.Qos(
		_prefetchCount,  // prefetch count
		_prefetchSize,   // prefetch size
		_prefetchGlobal, // global
	)
	if err != nil {
		ch.Close()
		c.channelPool.CheckConnect()
		return nil, errors.Wrap(err, "Error ch.Qos")
	}

	return ch, nil
}
