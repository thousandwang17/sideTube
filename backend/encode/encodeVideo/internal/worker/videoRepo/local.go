/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 16:29:37
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-22 14:18:30
 * @FilePath: /encode/encodeVideo/handler.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package videoRepo

import (
	"context"
	"encodeVideo/internal/worker"
	"encoding/json"
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
	ErrDetct        = errors.New("ffprobe error")
	ErrFileNotFound = errors.New("file not found")
	ErrFileFormat   = errors.New("ffempg failed tp get video format")
	ErrEncode       = errors.New("ffempg failed to encode")
	ErrVideoFormat  = errors.New("video format not reach min requirement")
	ErrAudioFormat  = errors.New("audio format not reach min requirement")
)

var (
	audioFormats = map[int]string{
		44100:  "44k",
		96000:  "96k",
		192000: "192k",
	}

	VideoFormat = map[int]struct{}{
		480:  struct{}{},
		720:  struct{}{},
		1080: struct{}{},
	}

	VideoFPSFormat = map[float64]int{
		29.97: 30,
		30:    30,
		60:    60,
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

func (b local) EncodeVideo(c context.Context, mission worker.EncodeVideoMission) (outputFile string, err error) {

	videofile := fmt.Sprintf("%s%s.mp4", b.path, mission.VideoId)
	if exist, err := b.checkFileStat(videofile); !exist || err != nil {
		return "", ErrFileNotFound
	}

	if mission.VideoFormat.Height > 0 {
		outputFile, err = b.encodeVideoStream(videofile, mission)
		if err != nil {
			return "", err
		}
	} else {
		outputFile, err = b.encodeAudioStream(videofile, mission)
		if err != nil {
			return "", err
		}
	}
	return outputFile, nil
}

//todo dected file alerdy encode by ffprobe
func (b local) encodeVideoStream(inputFile string, format worker.EncodeVideoMission) (outputFile string, err error) {

	fps, ok := VideoFPSFormat[format.VideoFormat.Fps]
	if !ok {
		log.Println("FPS Format wrong", format.VideoFormat.Fps)
		return "", ErrVideoFormat
	}

	_, ok = VideoFormat[format.VideoFormat.Height]
	if !ok {
		log.Println("Height Format wrong", format.VideoFormat.Height)
		return "", ErrVideoFormat
	}

	outputFile = fmt.Sprintf("%s%s.%d.%v.webm",
		b.path,
		format.VideoId,
		format.VideoFormat.Height,
		fps)

	// Run ffprobe tey  to get video information , if can get message that mean video already encode
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-select_streams", "v:0", outputFile)
	if out, err := cmd.Output(); err == nil {
		var ffprobeOutput VideoStream
		err = json.Unmarshal(out, &ffprobeOutput)
		if err != nil && len(ffprobeOutput.Streams) > 0 {
			return outputFile, nil
		}
	}

	// FFmpeg command to convert the file
	// First pass
	cmd = exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-c:v", "libvpx-vp9", "-b:v", "0", "-crf", "30",
		"-vf", fmt.Sprintf("scale=-2:%d,fps=fps=%v", format.VideoFormat.Height, format.VideoFormat.Fps),
		"-tile-columns", "4", "-frame-parallel", "1",
		"-pass", "1",
		"-an", "-f", "null",
		"-dash", "1",
		"/dev/null",
	)

	if err := cmd.Run(); err != nil {
		log.Println(format.VideoId, format.MissionID, err, cmd.String())
		return "", ErrEncode
	}
	log.Println(format.VideoId, format.MissionID, "pass 1 end")

	// Second pass
	cmd = exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-c:v", "libvpx-vp9", "-b:v", "0", "-crf", "30",
		"-vf", fmt.Sprintf("scale=-2:%d,fps=fps=%v", format.VideoFormat.Height, format.VideoFormat.Fps),
		"-tile-columns", "4", "-frame-parallel", "1",
		"-pass", "2",
		"-an",
		"-dash", "1",
		"-y",
		outputFile,
	)
	if err := cmd.Run(); err != nil {
		log.Println(format.VideoId, format.MissionID, err)
		return "", ErrEncode
	}

	return outputFile, nil

}

func (b local) encodeAudioStream(inputFile string, format worker.EncodeVideoMission) (outputFile string, err error) {
	outputFile = fmt.Sprintf("%s%s.%d.webm",
		b.path,
		format.VideoId,
		format.AudioFormat.Hz)

	// Run ffprobe to get video information , if can get message that mean video already encode
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-select_streams", "a:0", outputFile)
	if out, err := cmd.Output(); err == nil {
		var ffprobeOutput VideoStream
		err = json.Unmarshal(out, &ffprobeOutput)
		if err != nil && len(ffprobeOutput.Streams) > 0 {
			return outputFile, nil
		}
	}

	bitrate, ok := audioFormats[format.AudioFormat.Hz]
	if !ok {
		return "", ErrAudioFormat
	}

	cmd = exec.Command(
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
