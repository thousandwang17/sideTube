/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 18:45:31
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-10 21:00:02
 * @FilePath: /detectVideo/internal/worker/queueRepo/rabbitmq.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package queueRepo

import (
	"context"
	"detectVideo/internal/common/rabbitmq"
	"detectVideo/internal/worker"

	"github.com/rabbitmq/amqp091-go"
)

type queue struct {
	client *rabbitmq.RabbitClient
}

func NewRabbitmq(q *rabbitmq.RabbitClient) worker.Queue {
	return queue{q}
}

func (q queue) Consume() (<-chan amqp091.Delivery, error) {
	return q.client.Consume()
}

func (q queue) ReQueue(ctx context.Context, body []byte, delay int64) error {
	return q.client.PublishDectedVideo(ctx, body, amqp091.Table{"x-delay": delay})
}

func (q queue) PublishEncode(ctx context.Context, videomissions, audioMission [][]byte) error {
	return q.client.PublishEncode(ctx, videomissions, audioMission)
}
