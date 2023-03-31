/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 16:29:37
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-11 17:29:15
 * @FilePath: /encode/encodeVideo/handler.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package videoRepo

import (
	"context"
	"encodeVideo/internal/worker"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type local struct {
	path string
}

var (
	ErrDetct              = errors.New("ffprobe error")
	ErrFileNotFound       = errors.New("file not found")
	ErrFileFormat         = errors.New("ffempg failed tp get video format")
	ErrEncode             = errors.New("ffempg failed to encode")
	ErrAudioFormat        = errors.New("audio format not reach min requirement")
	ErrAudioMissionFormat = errors.New("Mission is not audio type ")
)

var (
	audioFormats = map[int]string{
		44100:  "44k",
		96000:  "96k",
		192000: "192k",
	}
)

type AudioStream struct {
	Streams []struct {
		CodecName  string `json:"codec_name"`
		SampleFmt  string `json:"sample_fmt"`
		SampleRate string `json:"sample_rate"`
	} `json:"streams"`
}

/**
 * @description:
 * @return {*}
 */
func NewLoacl(path string) worker.VideoRepository {
	return local{
		path: path,
	}
}

func (b local) EncodeAudio(c context.Context, mission worker.EncodeVideoMission) (outputFile string, err error) {

	videofile := fmt.Sprintf("%s%s.mp4", b.path, mission.VideoId)
	if exist, err := b.checkFileStat(videofile); !exist || err != nil {
		return "", ErrFileNotFound
	}

	outputFile, err = b.encodeAudioStream(videofile, mission)
	if err != nil {
		return "", err
	}

	return outputFile, nil
}

func (b local) encodeAudioStream(inputFile string, format worker.EncodeVideoMission) (outputFile string, err error) {
	outputFile = fmt.Sprintf("%s%s.%d.webm",
		b.path,
		format.VideoId,
		format.AudioFormat.Hz)

	bitrate, ok := audioFormats[format.AudioFormat.Hz]
	if !ok {
		return "", ErrAudioFormat
	}

	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-c:a", "libopus",
		"-b:a", bitrate,
		"-vn",
		"-f", "webm",
		"-dash", "1",
		"-y",
		outputFile,
	)
	if err := cmd.Run(); err != nil {
		log.Println(format.VideoId, format.MissionID, err, cmd.String())
		return "", ErrEncode
	}

	return outputFile, nil
}

func (b local) checkFileStat(pathTofile string) (bool, error) {
	if _, err := os.Stat(pathTofile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (b local) videoPath(videoID string) string {
	return fmt.Sprintf("%s%s.mp4", b.path, videoID)
}
