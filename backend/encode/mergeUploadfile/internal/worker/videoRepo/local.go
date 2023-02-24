/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 16:29:37
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-18 15:53:21
 * @FilePath: /encode/encodeVideo/handler.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package videoRepo

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"mergeUploadFile/internal/worker"
	"os"
)

type local struct {
	path string
}

var (
	ErrReachPackgetLimit = errors.New("count of video parts is reaching max limit")
)

/**
 * @description:
 * @return {*}
 */
func NewLoacl(path string) worker.VideoRepository {
	return local{
		path: path,
	}
}

func (b local) getCount(key string) (count uint32, err error) {
	file, err := os.OpenFile(fmt.Sprintf("%s%s.lock", b.path, key), os.O_RDWR, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	p := make([]byte, 4)
	rlen, rerr := file.ReadAt(p, 0)
	if rerr != nil {
		return 0, rerr
	}

	count = binary.LittleEndian.Uint32(p[:rlen])
	return count, nil
}

/**
 * @description:
 * @param {context.Context} c
 * @param {string} id
 * @return {*}
 */
func (b local) CompleteMultipartUpload(c context.Context, id string) error {
	mergeFileName := fmt.Sprintf("%s%s.mp4", b.path, id)

	count, err := b.getCount(id)
	if err != nil {
		// if is not exist, it mean its alerdy merge.
		if exist, errLock := b.checkFileStat(id + ".lock"); errors.Is(err, os.ErrNotExist) && exist {
			return nil
		} else if errLock != nil {
			return errLock
		}
		return err
	}

	fileName := fmt.Sprintf("%s%s.lock", b.path, id)

	file, err := os.Create(mergeFileName)
	if err != nil {
		log.Println(fileName, " err:", err)
		return err
	}

	defer file.Close()

	for part_id := uint32(1); part_id <= count; part_id++ {
		part, err := os.Open(fmt.Sprintf("%s%s_%d.tmp", b.path, id, part_id))
		if err != nil {
			log.Printf("%s part: %d , copy err %v ", fileName, part_id, err)
			return err
		}
		io.Copy(file, part)
		part.Close()
	}

	if err := b.abortUpload(c, id); err != nil {
		return err
	}
	log.Printf("%v.mp4  merge done!!", id)
	return nil
}

func (b local) abortUpload(c context.Context, id string) error {

	_, err := b.getCount(id)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fmt.Sprintf("%s%s.lock", b.path, id), os.O_RDWR, 0644)

	if err != nil {
		return err
	}
	defer file.Close()

	p := make([]byte, 4)

	var offset int64 = 4
	for {
		rlen, rerr := file.ReadAt(p, offset)
		if rerr != nil && rerr != io.EOF {
			return rerr
		}

		if rlen == 0 {
			break
		}

		partId := binary.LittleEndian.Uint32(p[:rlen])
		os.Remove(fmt.Sprintf("%s%s_%d.tmp", b.path, id, partId))
		offset += 4
	}

	os.Remove(fmt.Sprintf("%s%s.lock", b.path, id))
	return nil
}

func (b local) checkFileStat(fileName string) (bool, error) {
	if _, err := os.Stat(fmt.Sprintf("%s%s", b.path, fileName)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
