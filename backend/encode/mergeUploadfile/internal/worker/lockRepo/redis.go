/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 20:49:23
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-18 22:10:21
 * @FilePath: /mergeUploadfile/lockRepo/redis.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package lockRepo

import (
	"context"
	"mergeUploadFile/internal/worker"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisLocker struct {
	client *redis.Client
}

func NewRedis(c *redis.Client) worker.LockSystem {
	return redisLocker{c}
}

func (r redisLocker) Lock(ctx context.Context, videoID string, ttl time.Duration) error {
	res := r.client.SetNX(ctx, videoID, 1, ttl)
	return res.Err()
}

func (r redisLocker) UnLock(ctx context.Context, videoID string) error {
	res := r.client.Del(ctx, videoID)
	return res.Err()
}
