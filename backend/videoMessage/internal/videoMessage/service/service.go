/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-10 20:28:34
 * @FilePath: /videoMessage/internal/videoMessage/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"errors"
	"sideTube/videoMessage/internal/videoMessage"
)

var (
	ErrVaild = errors.New("args vaild failed")
)

const (
	private = iota
	publish
)

type service struct {
	metaRepo videoMessage.MetaRepository
}

type VideoMessageCommend interface {
	MessageList(c context.Context, videoId string, start, length int64) (messages []videoMessage.VideoMessageMeta, err error)
	Message(c context.Context, videoId, message string) (messageId string, err error)
	EditMessage(c context.Context, messageID, newMessage string) error
	DeleteMessage(c context.Context, messageID string) error

	// respond msg
	ReplyList(c context.Context, messageID string, start, length int64) (messages []videoMessage.VideoMessageReplyMeta, err error)
	Reply(c context.Context, messageID, message string) (messageId string, err error)
	EditReply(c context.Context, messageID, newMessage string) error
	DeleteReply(c context.Context, messageID string) error
}

func NewvideoMessageCommend(db videoMessage.MetaRepository) VideoMessageCommend {
	return service{
		metaRepo: db,
	}
}

func (v service) MessageList(c context.Context, videoId string, start, length int64) (messages []videoMessage.VideoMessageMeta, err error) {

	data, err := v.metaRepo.MessageList(c, videoId, start, length)
	if err != nil {
		return []videoMessage.VideoMessageMeta{}, err
	}

	return data, nil
}

func (v service) Message(c context.Context, videoID, message string) (messageId string, err error) {
	userId := c.Value("uid").(string)
	userName := c.Value("userName").(string)

	messageId, err = v.metaRepo.Message(c, userId, userName, videoID, message)

	if err != nil {
		return "", err
	}
	return messageId, nil
}

func (v service) EditMessage(c context.Context, messageId, newMessage string) error {
	userId := c.Value("uid").(string)
	err := v.metaRepo.EditMessage(c, userId, messageId, newMessage)
	if err != nil {
		return err
	}
	return nil
}

func (v service) DeleteMessage(c context.Context, messageId string) error {

	userId := c.Value("uid").(string)
	err := v.metaRepo.DeleteMessage(c, userId, messageId)

	if err != nil {
		return err
	}
	return nil
}

// respond msg

func (v service) ReplyList(c context.Context, messageID string, start, length int64) (messages []videoMessage.VideoMessageReplyMeta, err error) {
	data, err := v.metaRepo.ReplyList(c, messageID, start, length)
	if err != nil {
		return []videoMessage.VideoMessageReplyMeta{}, err
	}

	return data, nil
}

func (v service) Reply(c context.Context, messageId, message string) (ReplyId string, err error) {

	userId := c.Value("uid").(string)
	userName := c.Value("userName").(string)
	ReplyId, err = v.metaRepo.Reply(c, userId, userName, messageId, message)

	if err != nil {
		return "", err
	}
	return ReplyId, nil
}

func (v service) EditReply(c context.Context, messageId, newMessage string) error {
	userId := c.Value("uid").(string)
	err := v.metaRepo.EditReply(c, userId, messageId, newMessage)
	if err != nil {
		return err
	}
	return nil
}

func (v service) DeleteReply(c context.Context, messageId string) error {

	userId := c.Value("uid").(string)
	err := v.metaRepo.DeleteReply(c, userId, messageId)

	if err != nil {
		return err
	}
	return nil
}
