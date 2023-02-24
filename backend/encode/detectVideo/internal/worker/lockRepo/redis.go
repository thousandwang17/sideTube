/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 20:49:23
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-20 21:46:10
 * @FilePath: /detectVideo/lockRepo/redis.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package lockRepo

import (
	"context"
	"detectVideo/internal/worker"
	"errors"
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

func (r redisLocker) Lock(ctx context.Context, videoID string, ttl time.Duration) error {
	res := r.client.SetNX(ctx, videoID+"_dected", 1, ttl)
	return res.Err()
}

func (r redisLocker) UnLock(ctx context.Context, videoID string) error {
	res := r.client.Del(ctx, videoID+"_dected")
	return res.Err()
}

func (r redisLocker) SetMissionMap(ctx context.Context, videoID string, length int) error {

	err := r.setBitsTo1(ctx, videoID+"_encode", length)
	return err
}

func (r redisLocker) setBitsTo1(ctx context.Context, key string, end int) error {

	if end <= 0 {
		return ErrMinSubmission
	}

	if end > 100 {
		return ErrReachMaxSubmission
	}

	script := `
        local key = KEYS[1]
   		local stop = tonumber(ARGV[1])
        
        for i=0,stop do
            redis.call('SETBIT', key, i, 1)
        end
		redis.call('EXPIRE', 'key', 3600* 24 * 3 ) 

        return "OK"
    `
	args := []interface{}{end}
	_, err := r.client.Eval(ctx, script, []string{key}, args...).Result()
	return err
}
