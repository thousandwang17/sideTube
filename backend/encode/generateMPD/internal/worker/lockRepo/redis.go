/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 20:49:23
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-22 19:46:08
 * @FilePath: /generateMPD/lockRepo/redis.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package lockRepo

import (
	"context"
	"errors"
	"generateMPD/internal/worker"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrReachMaxSubmission = errors.New(" Reach max sub mission limit(100) ")
	ErrMinSubmission      = errors.New(" count of sub-mission must greater than 0 ")
)

type redisLocker struct {
	client *redis.Client
}

func NewRedis(c *redis.Client) worker.LockSystem {
	return redisLocker{c}
}

func (r redisLocker) Lock(ctx context.Context, videoID string, ttl time.Duration) (videoList, audioList []string, err error) {
	res := r.client.SetNX(ctx, videoID+"_MPD", 1, ttl)

	if res.Err() != nil {
		return nil, nil, res.Err()
	}

	resList := r.client.HGetAll(ctx, videoID+"_list")
	if resList.Err() != nil {
		return nil, nil, resList.Err()
	}

	values := resList.Val()

	for i := range values {
		count := strings.Count(values[i], ".")
		// count == 3 mean this path is video file name that format is  {vidoeID}.{px}.{fps}.webm
		if count == 3 {
			// "/var/videos/63f5b455acc4cb8a7796168d.720.30.webm"
			videoList = append(videoList, values[i])

			// count == 2 mean this path is audio file name that format is  {vidoeID}.{hz}.webm
		} else if count == 2 {
			//"/var/videos/63f5b455acc4cb8a7796168d.44100.webm"
			audioList = append(audioList, values[i])
		}
	}

	return videoList, audioList, nil
}

func (r redisLocker) UnLock(ctx context.Context, videoID string) error {
	script := `
	local lock_key = KEYS[1]
	local list_key = KEYS[2]

	redis.call('DEL', lock_key)
	redis.call('DEL', list_key)

	return 1
`
	args := []interface{}{}
	_, err := r.client.Eval(ctx, script, []string{videoID + "_MPD", videoID + "_list"}, args...).Result()
	if err != nil {
		log.Println("Accomplish ", err)
		return err
	}

	return nil
}
