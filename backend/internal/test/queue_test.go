package test

import (
	"testing"
	"time"

	"github.com/streadway/amqp"
)

// TestQueueConnection 测试RabbitMQ连接和基本操作
func TestQueueConnection(t *testing.T) {
	// 构建RabbitMQ连接URL
	url := "amqp://guest:guest@localhost:5672/"

	// 连接到RabbitMQ服务器
	conn, err := amqp.Dial(url)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	// 声明测试队列
	queueName := "test_queue"
	q, err := ch.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		t.Fatalf("Failed to declare queue: %v", err)
	}

	// 测试消息发送
	message := "test message"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// 测试消息接收
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		t.Fatalf("Failed to register a consumer: %v", err)
	}

	// 使用通道接收消息
	received := false
	go func() {
		for d := range msgs {
			if string(d.Body) == message {
				received = true
				return
			}
		}
	}()

	// 等待消息接收
	time.Sleep(2 * time.Second)
	if !received {
		t.Error("Message was not received")
	}

	// 清理测试队列
	_, err = ch.QueueDelete(
		queueName,
		false, // ifUnused
		false, // ifEmpty
		false, // noWait
	)
	if err != nil {
		t.Fatalf("Failed to delete test queue: %v", err)
	}
}