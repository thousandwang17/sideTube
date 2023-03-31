/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 20:49:23
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-11 19:41:36
 * @FilePath: /encodeVideo/lockRepo/redis.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package lockRepo

import (
	"context"
	"encodeVideo/internal/worker"
	"errors"
	"fmt"
	"log"
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

// lock and check the mission have been encoded
func (r redisLocker) Lock(ctx context.Context, videoID string, missionID, subMissionID int, ttl time.Duration) (done bool, err error) {
	return r.lock(ctx,
		r.lockKey(videoID, missionID, subMissionID),
		r.lockKey(videoID, missionID, subMissionID),
		missionID,
		ttl)
}

func (r redisLocker) UnLock(ctx context.Context, videoID string, missionID, subMissionID int) error {
	res := r.client.Del(ctx, r.lockKey(videoID, missionID, subMissionID))
	return res.Err()
}

func (r redisLocker) AccomplishbSubMission(ctx context.Context, mission worker.EncodeVideoMission, encodedFileName string) (alldone bool, err error) {
	return r.setBitsTo1(ctx,
		mission,
		encodedFileName,
	)
}

// use bitmap to recode each mission state , 1 = been encoded
func (r redisLocker) setBitsTo1(ctx context.Context, mission worker.EncodeVideoMission, encodedFileName string) (alldone bool, err error) {

	if mission.SubMissionID < 0 || mission.SubMissionID > 100 {
		return false, ErrMissionIndex
	}

	script := `
		local target_key = KEYS[1]
		local current_key = KEYS[2]
		local index = tonumber(ARGV[1])
		local stop = tonumber(ARGV[2])

		-- update target_key bitmap not exists then create it 
		if redis.call("exists", KEYS[1]) == 0 then
			for i=0,stop do
				redis.call('SETBIT', target_key, i, 1)
			end
	 	end
		

		-- update current bit 
		redis.call('SETBIT', current_key, index, 1)

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

	target_key := fmt.Sprintf("%s_%d_%s", mission.VideoId, mission.MissionID, "encode")
	current_key := fmt.Sprintf("%s_%d_%s_%s", mission.VideoId, mission.MissionID, "encode", "current")

	//TotalChunk need to del 1 , because bitmap start from 0
	args := []interface{}{mission.SubMissionID, mission.TotalChunk - 1}
	res, err := r.client.Eval(ctx, script, []string{target_key, current_key}, args...).Result()
	if _, ok := res.(int64); err != nil || !ok {
		log.Println("Accomplish ", err, ok)
		return false, err
	}

	return res.(int64) == 1, nil
}

func (r redisLocker) lock(ctx context.Context, lock_key, missions_key string, index int, ttl time.Duration) (done bool, err error) {

	expire := int64(ttl.Seconds())
	script := `
		local lock_key = KEYS[1]
        local missions_key = KEYS[2]
   		local index = tonumber(ARGV[1])
		local ttl = tonumber(ARGV[2])

		-- check current mission had benn Encoded
		if redis.call('GETBIT', missions_key, index) == 1 then
		return 2
		end

		-- set lock if err return 0
		if redis.call('SETNX', lock_key, 1) == 1 then 
		return redis.call('EXPIRE', lock_key, ttl)
		end

		return 0
    `
	args := []interface{}{index, expire}
	res, err := r.client.Eval(ctx, script, []string{lock_key, missions_key}, args...).Result()
	if _, ok := res.(int64); err != nil || !ok {
		log.Println("lock ", err, ok)
		return false, err
	}
	return res.(int64) == 2, nil
}

func (r redisLocker) lockKey(videoID string, missionID, subMissionID int) string {
	return fmt.Sprintf("%s_%d_%d_%s", videoID, missionID, subMissionID, "encode")
}

// 		local file_name = ARGV[2]
// -- insert the encoded file name to HSET
// redis.call('HSET', list_key, index, file_name)
