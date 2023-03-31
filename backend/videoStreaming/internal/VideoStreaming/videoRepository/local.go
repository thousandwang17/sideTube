/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-30 15:43:56
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-12 16:28:08
 * @FilePath: /videoUpload/internal/videoUpload/videoRepository/local.go
 * @Description: this repo just for local test
 */
package videoRepository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	v "sideTube/VideoStreaming/internal/VideoStreaming"
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

const VIDEO_MAX_BUFFER_SIZE = 1024 * 1024 * 5 //3MB
const AUDIO_MAX_BUFFER_SIZE = 1024 * 1024     //1MB

type loaclS3 struct {
	file          *os.File
	LimitedReader io.Reader
}

func NewLoaclS3(b *os.File, limit int64) io.ReadCloser {
	return loaclS3{b, io.LimitReader(b, limit)}
}

func (l loaclS3) Read(p []byte) (n int, err error) {
	return l.LimitedReader.Read(p)
}

func (l loaclS3) Close() error {
	return l.file.Close()
}

/**
 * @description:
 * @return {*}
 */
func NewLoacl(path string) v.VideoRepository {
	return local{
		path: path,
	}
}

func (b local) GetVideo(c context.Context, videoId string, start int64, end int64) (data io.ReadCloser, err error) {

	file, err := os.Open(fmt.Sprintf("%s%s", b.path, videoId))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrVideoNotExists
		}
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, ErrVideoNotExists
	}

	if end == -1 || end >= fileInfo.Size()-1 {
		end = fileInfo.Size() - 1
	}

	if start > end {
		return nil, ErrorFileRange
	}

	length := int(end - start)
	if length > VIDEO_MAX_BUFFER_SIZE {
		return nil, ErrVideoReachMaxBufferSzie
	}

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return NewLoaclS3(file, end-start+1), nil
}

func (b local) GetMpdFile(c context.Context, fileName string) (data io.ReadCloser, err error) {

	file, err := os.Open(fmt.Sprintf("%s%s", b.path, fileName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrVideoNotExists
		}
		return nil, err
	}

	return file, nil
}
