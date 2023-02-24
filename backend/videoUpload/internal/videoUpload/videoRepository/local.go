/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-30 15:43:56
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-17 23:56:09
 * @FilePath: /videoUpload/internal/videoUpload/videoRepository/local.go
 * @Description: this repo just for local test
 */
package videoRepository

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	v "sideTube/videoUpload/internal/videoUpload"
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
func NewLoacl(path string) v.VideoRepository {
	return local{
		path: path,
	}
}

/**
 * @description:
 * @param {context.Context} c
 * @return {*}
 */
func (b local) CreateMultipartUpload(c context.Context, id string) (v.VideoRepoMeta, error) {

	if exist, err := b.keyExists(id); err != nil {
		return v.VideoRepoMeta{}, err
	} else if exist {
		return v.VideoRepoMeta{}, errors.New("videoId is already exists")
	}

	file, err := os.Create(fmt.Sprintf("%s%s.lock", b.path, id))

	if err != nil {
		log.Println(err, 213)
		return v.VideoRepoMeta{}, err
	}
	defer file.Close()

	p := make([]byte, 4)
	binary.LittleEndian.PutUint32(p, 0)
	file.WriteAt(p, 0)
	return v.VideoRepoMeta{Id: id}, nil
}

/**
 * @description:
 * @param {context.Context} c
 * @param {string} id
 * @param {int64} partId
 * @param {io.ReadSeeker} file
 * @return {*}
 */
func (b local) UploadPart(c context.Context, v v.VideoRepoMeta, file io.ReadSeeker) error {

	if count, err := b.getCount(v.Id); err != nil {
		return err
	} else if count >= 1000 {
		return ErrReachPackgetLimit
	}

	tmpfile, err := os.Create(fmt.Sprintf("%s%s_%d.tmp", b.path, v.Id, v.PartID))

	if err != nil {
		return err
	}
	defer tmpfile.Close()

	p := make([]byte, 1024*1024)

	for {
		rlen, rerr := file.Read(p)

		if rerr != nil && rerr != io.EOF {
			return rerr
		}
		if rlen == 0 {
			break
		}

		_, werr := tmpfile.Write(p[:rlen])

		if werr == io.ErrShortWrite {
			return werr
		}
	}

	if err := b.addCount(v.Id, uint32(v.PartID)); err != nil {
		return err
	}

	return nil
}

/**
 * @description:
 * @param {context.Context} c
 * @param {string} id
 * @return {*}
 */
func (b local) CompleteMultipartUpload(c context.Context, id string) error {
	mergeFileName := fmt.Sprintf("%v.mp4", id)

	count, err := b.getCount(id)
	if err != nil {
		return err
	}

	go func() {
		fileName := fmt.Sprintf("%s%s.lock", b.path, id)

		file, err := os.Create(mergeFileName)

		if err != nil {
			log.Println(fileName, " err:", err)
			return
		}

		defer file.Close()

		for part_id := uint32(1); part_id <= count; part_id++ {
			part, err := os.Open(fmt.Sprintf("%s%s_%d.tmp", b.path, id, part_id))
			if err != nil {
				log.Printf("%s part: %d , copy err %v ", fileName, part_id, err)
				return
			}
			io.Copy(file, part)
			part.Close()
		}

		b.AbortUpload(context.Background(), id)
		log.Printf("%v.mp4  merge done!!", id)
	}()
	return nil
}

/**
 * @description:
 * @param {context.Context} c
 * @param {string} id
 * @return {*}
 */
func (b local) AbortUpload(c context.Context, id string) error {

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

/**
 * @description:
 * @param {string} key
 * @return {*}
 */
func (b local) keyExists(key string) (bool, error) {
	if _, err := os.Stat(fmt.Sprintf("%s%s.lock", b.path, key)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	if _, err := os.Stat(fmt.Sprintf("%s%s.mp4", b.path, key)); err == nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (b local) addCount(key string, partId uint32) error {

	fileName := fmt.Sprintf("%s%s.lock", b.path, key)

	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	p := make([]byte, 4)
	rlen, rerr := file.ReadAt(p, 0)
	if rerr != nil {
		return rerr
	}

	// read and write file  to add count at  pointer[0:4]
	count := binary.LittleEndian.Uint32(p[:rlen])
	if count >= 1000 {
		return ErrReachPackgetLimit
	}

	binary.LittleEndian.PutUint32(p, count+1)
	_, werr := file.WriteAt(p, 0) // Write at 0 beginning
	if werr != nil {
		return werr
	}

	// append part_id in to lock to record
	binary.LittleEndian.PutUint32(p, partId)
	seekEnd, err := file.Seek(0, io.SeekEnd)
	_, werr = file.WriteAt(p, seekEnd)
	if werr != nil {
		return werr
	}

	return nil
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
