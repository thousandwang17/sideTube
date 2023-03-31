/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 16:29:37
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-21 15:29:03
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
	ErrVideoFormat        = errors.New("video format not reach min requirement")
	ErrVideoMissionFormat = errors.New("Mission is not video type ")
)

var (
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

	videofile := fmt.Sprintf("%s%s_%d.mp4", b.path, mission.VideoId, mission.SubMissionID)
	if exist, err := b.checkFileStat(videofile); !exist || err != nil {
		return "", ErrFileNotFound
	}

	outputFile, err = b.encodeVideoStream(mission.VideoId, mission)
	if err != nil {
		return "", err
	}

	return outputFile, nil
}

//todo dected file alerdy encode by ffprobe
func (b local) encodeVideoStream(VideoId string, format worker.EncodeVideoMission) (outputFile string, err error) {

	_, ok := VideoFormat[format.VideoFormat.Height]
	if !ok {
		log.Println("Height Format wrong", format.VideoFormat.Height)
		return "", ErrVideoFormat
	}

	outputFile = fmt.Sprintf("%s%s_%d.%d.%v.webm",
		b.path,
		format.VideoId,
		format.SubMissionID,
		format.VideoFormat.Height,
		format.VideoFormat.Fps,
	)

	vf := fmt.Sprintf("scale=-2:%d,format=yuv420p,fps=fps=%v", format.VideoFormat.Height, format.VideoFormat.Fps)
	if format.VideoFormat.OriginFps {
		vf = fmt.Sprintf("scale=-2:%d,format=yuv420p", format.VideoFormat.Height)
	}

	// FFmpeg command to convert the file
	inputFileName := fmt.Sprintf("%s_%d.mp4", VideoId, format.SubMissionID)
	inputFile := b.videoPath(inputFileName)
	// First pass
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-c:v", "libvpx-vp9", "-b:v", "0", "-crf", "30",
		"-vf", vf,
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
	log.Println(format.VideoId, format.MissionID, " 2 pass 1 end")

	// Second pass
	cmd = exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-c:v", "libvpx-vp9", "-b:v", "0", "-crf", "30",
		"-vf", vf,
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
	log.Println(format.VideoId, format.MissionID, " 2 pass 2 end")

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

func (b local) videoPath(fileName string) string {
	return fmt.Sprintf("%s%s", b.path, fileName)
}
