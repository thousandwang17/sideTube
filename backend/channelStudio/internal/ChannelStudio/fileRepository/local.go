/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-30 15:43:56
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 21:10:43
 * @FilePath: /videoUpload/internal/videoUpload/videoRepository/local.go
 * @Description: this repo just for local test
 */
package fileRepository

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	v "sideTube/ChannelStudio/internal/ChannelStudio"

	"github.com/nfnt/resize"
)

type local struct {
	path string
}

var (
	ErrUpload    = errors.New("file server error")
	ErrEmptyData = errors.New("empty data")
)

/**
 * @description:
 * @return {*}
 */
func NewLoacl(path string) v.FileReop {
	return local{
		path: path,
	}
}

func (b local) SetVideoPicture(c context.Context, video_id string, extension string, data io.ReadSeekCloser) (imageName string, err error) {

	defer data.Close()
	fileName := fmt.Sprintf("%s.%s", video_id, "png")
	file, err := os.Create(fmt.Sprintf("%s%s", b.path, fileName))

	if err != nil {
		log.Println(video_id, "Create", err)
		return "", ErrUpload
	}

	defer func() {
		file.Close()
		if err != nil {
			os.Remove(b.path + file.Name())
		}
	}()

	// Decode the uploaded image
	img, _, err := image.Decode(data)
	if err != nil {
		log.Println(video_id, "Decode", err)
		return "", ErrUpload
	}

	// Resize the image to 720x404 using Lanczos filter
	resized := resize.Resize(720, 404, img, resize.Lanczos3)

	// Save the resized image to file
	err = png.Encode(file, resized)
	if err != nil {
		log.Println(video_id, "Encode", err)
		return "", ErrUpload
	}

	return fileName, nil
}
