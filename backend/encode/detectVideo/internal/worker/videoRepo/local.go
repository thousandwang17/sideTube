/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 16:29:37
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-21 23:19:45
 * @FilePath: /encode/encodeVideo/handler.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package videoRepo

import (
	"context"
	"detectVideo/internal/worker"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type local struct {
	path string
}

var (
	ErrDetct        = errors.New("ffprobe error")
	ErrFileNotFound = errors.New("file not found")
	ErrFileFormat   = errors.New("ffprobe failed tp get video format")
	ErrVideoFormat  = errors.New("video format not reach min requirement")
	ErrAudioFormat  = errors.New("audio format not reach min requirement")
)

var (
	videoFormats = [...]worker.VideoFormat{
		{Width: 1920, Height: 1080, Fps: 60},
		{Width: 1920, Height: 1080, Fps: 29.97},
		{Width: 1280, Height: 720, Fps: 60},
		{Width: 1280, Height: 720, Fps: 29.97},
		{Width: 854, Height: 480, Fps: 60},
		{Width: 854, Height: 480, Fps: 29.97},
	}

	audioFormats = [...]worker.AudioFormat{
		{Hz: 44100},
		{Hz: 96000},
		{Hz: 192000},
	}
)

type AudioStream struct {
	Streams []struct {
		CodecName  string `json:"codec_name"`
		SampleFmt  string `json:"sample_fmt"`
		SampleRate string `json:"sample_rate"`
	} `json:"streams"`
}

type VideoStream struct {
	Streams []struct {
		CodecType string `json:"codec_type"`
		Width     int    `json:"width,omitempty"`
		Height    int    `json:"height,omitempty"`
		Fps       string `json:"r_frame_rate,omitempty"`
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

func (b local) DetectVideo(c context.Context, videoID string) ([]worker.VideoFormat, []worker.AudioFormat, error) {

	videofile := fmt.Sprintf("%s%s.mp4", b.path, videoID)
	if exist, err := b.checkFileStat(videofile); !exist || err != nil {
		return nil, nil, ErrFileNotFound
	}

	videoMission, err := b.detectVideoStream(videofile)
	if err != nil {
		return nil, nil, err
	}

	audioMissopm, err := b.detectAudioStream(videofile)
	if err != nil {
		return nil, nil, err
	}

	return videoMission, audioMissopm, nil
}

func (b local) detectVideoStream(videofile string) ([]worker.VideoFormat, error) {

	// Run ffprobe to get video information
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-select_streams", "v:0", videofile)

	out, err := cmd.Output()
	if err != nil {
		log.Println("Error:", err)
		return nil, ErrFileFormat
	}

	var ffprobeOutput VideoStream
	err = json.Unmarshal(out, &ffprobeOutput)
	if err != nil {
		log.Println("Error:", err)
		return nil, ErrFileFormat
	}

	if len(ffprobeOutput.Streams) == 0 {
		return nil, ErrFileFormat
	}

	videoStream := ffprobeOutput.Streams[0]

	// get fps
	rFrameRate := videoStream.Fps
	parts := strings.Split(rFrameRate, "/")
	numerator, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		log.Println("strconv Error:", err)
		return nil, err
	}
	denominator, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		log.Println("strconv denominator Error:", err)
		return nil, err
	}

	fps := numerator / denominator

	missions := []worker.VideoFormat{}

	for i := range videoFormats {
		if videoStream.Height < videoFormats[i].Height {
			continue
		}

		if videoStream.Width < videoFormats[i].Width {
			continue
		}

		if fps < videoFormats[i].Fps {
			continue
		}

		missions = append(missions, videoFormats[i])
	}

	if len(missions) == 0 {
		return nil, ErrVideoFormat
	}

	return missions, nil
}

func (b local) detectAudioStream(videofile string) ([]worker.AudioFormat, error) {

	// Run ffprobe to get video information
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-select_streams", "a:0", videofile)

	out, err := cmd.Output()
	if err != nil {
		log.Println("Error:", err)
		return nil, ErrFileFormat
	}

	var ffprobeOutput AudioStream
	err = json.Unmarshal(out, &ffprobeOutput)
	if err != nil {
		log.Println("Error:", err)
		return nil, ErrFileFormat
	}

	if len(ffprobeOutput.Streams) == 0 {
		return nil, ErrFileFormat
	}

	audio := ffprobeOutput.Streams[0]

	SampleRate, err := strconv.Atoi(audio.SampleRate)
	if err != nil {
		log.Println("strconv SampleRate Error:", err)
		return nil, err
	}

	missions := []worker.AudioFormat{}

	for i := range audioFormats {
		if SampleRate < audioFormats[i].Hz {
			continue
		}

		missions = append(missions, audioFormats[i])
	}

	if len(missions) == 0 {
		return nil, ErrAudioFormat
	}

	return missions, nil
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
