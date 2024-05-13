package publisher

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/rabbitmq"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
)

const (
	_publishMandatory = false
	_publishImmediate = false

	_exchangeName    = "orders-exchange"
	_bindingKey      = "orders-routing-key"
	_messageTypeName = "ordered"
)

type Publisher struct {
	exchangeName, bindingKey string
	messageTypeName          string
	channelPool              *rabbitmq.ChannelPool
}

var EventPublisherSet = wire.NewSet(NewPublisher)

func NewPublisher(cp *rabbitmq.ChannelPool) (Publisher, error) {
	pub := Publisher{
		exchangeName:    _exchangeName,
		bindingKey:      _bindingKey,
		messageTypeName: _messageTypeName,
		channelPool:     cp,
	}

	return pub, nil
}

func (p *Publisher) Configure(opts ...Option) EventPublisher {
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Publisher) PublishEvents(ctx context.Context, events []any) error {
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			return errors.Wrap(err, "publisher-json.Marshal")
		}

		err = p.Publish(ctx, b, "text/plain")
		if err != nil {
			return errors.Wrap(err, "publisher-pub.Publish")
		}
	}

	return nil
}

// Publish message.
func (p *Publisher) Publish(ctx context.Context, body []byte, contentType string) error {
	ch, err := p.channelPool.GetChannel()
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}

	slog.Info("publish message", "exchange", p.exchangeName, "routing_key", p.bindingKey)

	if err := ch.PublishWithContext(
		ctx,
		p.exchangeName,
		p.bindingKey,
		_publishMandatory,
		_publishImmediate,
		amqp.Publishing{
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			Timestamp:    time.Now(),
			Body:         body,
			Type:         p.messageTypeName,
		},
	); err != nil {
		ch.Close()
		p.channelPool.CheckConnect()
		return errors.Wrap(err, "ch.Publish")
	}
	//归还channel
	p.channelPool.ReturnChannel(ch)
	return nil
}
