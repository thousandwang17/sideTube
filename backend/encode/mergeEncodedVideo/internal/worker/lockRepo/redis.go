/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 20:49:23
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-11 14:52:42
 * @FilePath: /mergeUploadfile/lockRepo/redis.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package lockRepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mergeUploadFile/internal/worker"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrMissionIndex = errors.New(" mission index should between 0 ~ 100")
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

func (r redisLocker) AccomplishbMission(ctx context.Context, mission worker.Mission, encodedFileName string) (alldone bool, err error) {
	return r.setBitsTo1(ctx,
		fmt.Sprintf("%s_%s", mission.VideoId, "encode"),
		fmt.Sprintf("%s_%s_%s", mission.VideoId, "encode", "current"),
		fmt.Sprintf("%s_%s", mission.VideoId, "list"),
		mission,
		encodedFileName,
	)
}

// use bitmap to recode each mission state , 1 = been encoded
func (r redisLocker) setBitsTo1(ctx context.Context, target_key, current_key, list_key string, mission worker.Mission, encodedFileName string) (alldone bool, err error) {

	if mission.MissionID < 0 || mission.MissionID > 100 {
		return false, ErrMissionIndex
	}

	script := `
		local target_key = KEYS[1]
		local current_key = KEYS[2]
		local list_key = KEYS[3]
		local index = tonumber(ARGV[1])
		local file_name = ARGV[2]

		-- update current bit 
		redis.call('SETBIT', current_key, index, 1)
	
		-- insert the encoded file name to HSET
		redis.call('HSET', list_key, index, file_name)
		
		-- Check that the bitmaps are the same size
		local size1 = redis.call('bitcount', target_key)
		local size2 = redis.call('bitcount', current_key)
		if size1 ~= size2 then
			return 0
		end

		if size1 == 0 or size2 == 0 then
			return 0
		end

		-- Iterate over the bits in each bitmap and compare them
		local temp_map = target_key .. '_temp'
		local result = redis.call('BITOP', 'XOR', temp_map, target_key, current_key)
		local difference = redis.call('BITCOUNT', temp_map)
		
		-- Clean up the temporary map
		redis.call('DEL', temp_map)

		if difference ~= 0 then
			return 0
		end

		-- If we got here, the bitmaps are equal , and del target_key and current_key
		redis.call('DEL', target_key)
		redis.call('DEL', current_key)
	
		return 1
	`
	args := []interface{}{mission.MissionID, encodedFileName}
	res, err := r.client.Eval(ctx, script, []string{target_key, current_key, list_key}, args...).Result()
	if _, ok := res.(int64); err != nil || !ok {
		log.Println("Accomplish ", err, ok)
		return false, err
	}

	return res.(int64) == 1, nil
}
