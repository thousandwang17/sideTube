/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-30 15:43:56
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-02 16:16:07
 * @FilePath: /videoUpload/internal/videoUpload/videoRepository/local.go
 * @Description: this repo just for local test
 */
package fileRepository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	v "sideTube/picture/internal/picture"
)

type local struct {
	path string
}

var (
	ErrVideoNotExists          = errors.New("video is not exists")
	ErrorFileRange             = errors.New("invalid file Range ")
	ErrVideoReachMaxBufferSzie = errors.New("video length over max buffer limit")
	ErrAudioReachMaxBufferSzie = errors.New("audio length over max buffer limit")
	ErrFileTypeRequest         = errors.New("only can request .mpd file")
)

const VIDEO_MAX_BUFFER_SIZE = 1024 * 1024 * 3 //3MB
const AUDIO_MAX_BUFFER_SIZE = 1024 * 1024     //1MB

type loaclS3 struct {
	file          *os.File
	LimitedReader io.Reader
}

/**
 * @description:
 * @return {*}
 */
func NewLoacl(path string) v.PictureRepository {
	return local{
		path: path,
	}
}

func (b local) GetPicture(c context.Context, fileName string) (data io.ReadCloser, err error) {

	file, err := os.Open(fmt.Sprintf("%s%s", b.path, fileName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrVideoNotExists
		}
		return nil, err
	}

	return file, nil
}
