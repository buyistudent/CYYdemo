package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
	"sync"
)

// ChannelPool amqp 通道池，提供通道复用
type ChannelPool struct {
	url      RabbitMQConnStr
	mutex    sync.Mutex
	conn     *amqp.Connection
	channels []*amqp.Channel
}

// InitPool 初始化通道池
func (cp *ChannelPool) InitPool(url RabbitMQConnStr) {
	cp.url = url
	cp.channels = make([]*amqp.Channel, 0, 1000)
}

func (cp *ChannelPool) GetChannel() (*amqp.Channel, error) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	size := len(cp.channels)
	if size > 0 {
		ch := cp.channels[size-1]
		cp.channels[size-1] = nil
		cp.channels = cp.channels[:size-1]
		return ch, nil
	}
	slog.Info("RabbitMQ尝试连接，连接url：", cp.url)
	if cp.conn == nil {
		conn, err := amqp.Dial(string(cp.url))
		if err != nil {
			slog.Error("RabbitMQ尝试连接失败，", err.Error())
			return nil, err
		}
		cp.conn = conn
	}
	ch, err := cp.conn.Channel()
	if err != nil {
		cp.conn.Close()
		cp.conn = nil
		return nil, err
	}
	return ch, nil
}

func (cp *ChannelPool) ReturnChannel(ch *amqp.Channel) {
	cp.mutex.Lock()
	cp.channels = append(cp.channels, ch)
	cp.mutex.Unlock()
}

func (cp *ChannelPool) CheckConnect() {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	if cp.conn == nil {
		return
	}
	ch, err := cp.conn.Channel()
	if err != nil {
		for i := 0; i < len(cp.channels); i++ {
			cp.channels[i].Close()
			cp.channels[i] = nil
		}
		cp.channels = cp.channels[0:0]
		cp.conn.Close()
		cp.conn = nil
		return
	}
	cp.channels = append(cp.channels, ch)
}

// Publish 发布消息
func (cp *ChannelPool) Publish(exchange, key string, mandatory, immediate, reliable bool, msg amqp.Publishing) error {
	ch, err := cp.GetChannel()
	if err != nil {
		return err
	}
	if reliable {
		if err := ch.Confirm(false); err != nil {
			slog.Error("Channel could not be put into confirm mode: %s", err)
		}

		confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

		err = ch.Publish(exchange, key, mandatory, immediate, msg)
		if err != nil {
			ch.Close()
			cp.CheckConnect()
			return err
		}
		defer cp.ConfirmOne(confirms)

		return nil
	} else {
		err = ch.Publish(exchange, key, mandatory, immediate, msg)
		if err != nil {
			ch.Close()
			cp.CheckConnect()
			return err
		}
	}

	cp.ReturnChannel(ch)
	return nil
}

func (cp *ChannelPool) ConfirmOne(confirms <-chan amqp.Confirmation) error {
	slog.Info("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		slog.Info("confirmed delivery with delivery tag:", confirmed.DeliveryTag)
		return nil
	} else {
		slog.Error("failed delivery of delivery tag:", confirmed.DeliveryTag)
		return fmt.Errorf("failed delivery of delivery tag:", confirmed.DeliveryTag)
	}
}
