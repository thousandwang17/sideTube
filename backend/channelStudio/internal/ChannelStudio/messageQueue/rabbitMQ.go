/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-14 16:19:41
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-16 20:58:30
 * @FilePath: /channelStudio/internal/channelStudio/messageQueue/rabbitMQ.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package messageQueue

import (
	"context"
	"encoding/json"
	"sideTube/ChannelStudio/internal/ChannelStudio"
	"sideTube/ChannelStudio/internal/common/rabbitmq"
	"time"
)

type rabbitRepo struct {
	client *rabbitmq.RabbitClient
}

type Message struct {
	VideoId string    `json:"video_id"`
	UserId  string    `json:"user_id"`
	Time    time.Time `json:"time"`
}

func NewMessageRepo(r *rabbitmq.RabbitClient) ChannelStudio.MessagQqueue {
	return &rabbitRepo{
		client: r,
	}
}

func (r rabbitRepo) NotifySearchEngine(ctx context.Context, videoId, userId string) (err error) {
	jsonBytes, err := json.Marshal(Message{videoId, userId, time.Now()})
	if err != nil {
		return err
	}
	return r.client.PublishSearchEngine(ctx, jsonBytes)
}
