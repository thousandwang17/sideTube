/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 16:29:37
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-28 15:08:26
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
	"math"
	"os"
	"os/exec"
	"path/filepath"
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
	ErrChunckVideo  = errors.New("Falied to chunk video")
)

var (
	videoFormats = [...]worker.VideoFormat{
		{Width: 1920, Height: 1080},
		{Width: 1280, Height: 720},
		{Width: 854, Height: 480},
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
		sampleRate int
	} `json:"streams"`
}

type VideoStream struct {
	Streams []struct {
		CodecName string `json:"codec_name"`
		CodecType string `json:"codec_type"`
		BitRate   string `json:"bit_rate,omitempty"`
		bitRate   float64
		Fps       string `json:"r_frame_rate,omitempty"`
		Duration  string `json:"duration,omitempty"`
		fps       float64
		Width     int `json:"width,omitempty"`
		Height    int `json:"height,omitempty"`
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

func (b local) DetectVideo(c context.Context, videoID string) ([]worker.VideoFormat, []worker.AudioFormat, int, error) {

	videofile := b.videoPath(videoID)
	if exist, err := b.checkFileStat(videofile); !exist || err != nil {
		return nil, nil, 0, ErrFileNotFound
	}

	// ---------- audio -----------
	// if len of audioStream' stream equal 0 will return error
	audioStream, err := b.detectAudioStream(videofile)
	if err != nil {
		return nil, nil, 0, err
	}

	audioMissopm, err := b.generateAudioEncodeMissoion(videofile, audioStream)
	if err != nil {
		return nil, nil, 0, err
	}

	// ---------- video -----------
	// if len of videoStream' stream equal 0 will return error
	videoStream, err := b.detectVideoStream(videofile)
	if err != nil {
		return nil, nil, 0, err
	}

	count, err := b.chunkVideoByBitRate(videoID, videoStream)
	if err != nil {
		return nil, nil, 0, err
	}

	videoMission, err := b.generateVideoEncodeMissoion(videofile, videoStream, count)
	if err != nil {
		return nil, nil, 0, err
	}

	return videoMission, audioMissopm, count, nil
}

func (b local) detectVideoStream(videofile string) (VideoStream, error) {

	// Run ffprobe to get video information
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-select_streams", "v:0", videofile)

	var ffprobeOutput VideoStream
	out, err := cmd.Output()
	if err != nil {
		log.Println("Error:", err)
		return ffprobeOutput, ErrFileFormat
	}

	err = json.Unmarshal(out, &ffprobeOutput)
	if err != nil {
		log.Println("Error:", err)
		return ffprobeOutput, ErrFileFormat
	}

	if len(ffprobeOutput.Streams) == 0 {
		return ffprobeOutput, ErrFileFormat
	}

	videoStream := ffprobeOutput.Streams[0]

	if videoStream.CodecName != "h264" {
		log.Println(videofile, " wrong code type", videoStream.CodecName)
		return ffprobeOutput, ErrVideoFormat
	}

	// get fps
	rFrameRate := videoStream.Fps
	parts := strings.Split(rFrameRate, "/")

	numerator, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		log.Println("strconv Error:", err)
		return ffprobeOutput, err
	}

	denominator, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		log.Println("strconv denominator Error:", err)
		return ffprobeOutput, err
	}

	fpsFloat := numerator / denominator
	fps := math.Floor(fpsFloat + 0.5)
	log.Println(videofile, videoStream.Height, videoStream.Width, fps)

	ffprobeOutput.Streams[0].fps = fps

	// get bitrate
	bitrateString := videoStream.BitRate
	bitrate, err := strconv.ParseFloat(bitrateString, 64)
	if err != nil {
		log.Println("strconv bitrate Error:", err)
		return ffprobeOutput, err
	}

	if bitrate <= 0 {
		log.Println(" bitrate cant not be 0 ")
		return ffprobeOutput, err
	}

	ffprobeOutput.Streams[0].bitRate = bitrate

	log.Println(ffprobeOutput)
	return ffprobeOutput, nil

}

func (b local) chunkVideoByBitRate(videoid string, info VideoStream) (int, error) {
	stream := info.Streams[0]
	durationfloat, err := strconv.ParseFloat(stream.Duration, 64)
	if err != nil {
		log.Println("durationFloat Error:", err)
		return 0, err
	}
	durationInt := int(math.Ceil(durationfloat))

	if durationInt <= 0 {
		log.Println("durationInt can not be 0")
		return 0, errors.New("durationInt can not be 0")
	}

	encodeSize := os.Getenv("ENCODE_VIDEO_SIZE")

	encodeSizeFloat, err := strconv.ParseFloat(encodeSize, 64)
	if err != nil {
		log.Println("encodeSizeInt Error:", err)
		return 0, err
	}

	if encodeSizeFloat < 10000000000000 {
		encodeSizeFloat = 10000000000000
	}

	countFloat := stream.bitRate * float64(stream.Width) * float64(stream.Height) * durationfloat / encodeSizeFloat
	// WARNING : the number of chunk should detect by runing `ls xxx_*.mp4`
	// `countInt` can not be return directly ,because  the `count` not always equal ffmpeg segment number
	countInt := int(math.Ceil(countFloat))

	if countInt <= 0 {
		fmt.Println("chunkVideoByBitRate countInt is illegal:", countInt, stream.bitRate, durationfloat)
		return 0, ErrChunckVideo
	}

	log.Println(countInt, stream.bitRate, durationfloat, encodeSizeFloat)

	segment_time := 60 * 60 * 24

	sec := durationInt / countInt
	if sec > 0 {
		segment_time = sec
	} else {
		fmt.Println("chunkVideoByBitRate sec is illegal:", sec, durationInt, countInt)
		return 0, ErrChunckVideo
	}

	log.Printf("%f * %f / %f=  %d", stream.bitRate, durationfloat, encodeSizeFloat, segment_time)

	cmd := exec.Command("ffmpeg",
		"-i", b.videoPath(videoid),
		"-vcodec", "copy",
		"-segment_time", fmt.Sprint(segment_time),
		"-f", "segment",
		"-reset_timestamps", "1",
		"-avoid_negative_ts", "1",
		b.path+videoid+"_%d.mp4")

	log.Println(cmd.String())
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running ffmpeg command:", err)
		return 0, err
	}

	segments, err := b.countFilesWithPrefix(videoid)
	if err != nil {
		log.Println("Error running ls command:", err)
		return 0, ErrChunckVideo
	}

	if segments == 0 {
		log.Println("segments is illegal:", err)
		return 0, ErrChunckVideo
	}

	return segments, nil
}

func (b local) generateVideoEncodeMissoion(videofile string, VideoStream VideoStream, numOfChunk int) ([]worker.VideoFormat, error) {
	missions := []worker.VideoFormat{}

	if len(VideoStream.Streams) == 0 {
		return nil, ErrVideoFormat
	}

	videoStream := VideoStream.Streams[0]

	for i := range videoFormats {
		if videoStream.Height < videoFormats[i].Height {
			continue
		}

		if videoStream.Width < videoFormats[i].Width {
			continue
		}

		if videoStream.fps == 60 {
			missions = append(missions, worker.VideoFormat{
				Fps:    30,
				Height: videoFormats[i].Height,
				Width:  videoFormats[i].Width,
			})
		}

		missions = append(missions, worker.VideoFormat{
			Fps:       videoStream.fps,
			Height:    videoFormats[i].Height,
			Width:     videoFormats[i].Width,
			OriginFps: true,
		})

	}

	if len(missions) == 0 {
		return nil, ErrVideoFormat
	}

	return missions, nil
}

func (b local) detectAudioStream(videofile string) (AudioStream, error) {

	// Run ffprobe to get video information
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-select_streams", "a:0", videofile)

	var ffprobeOutput AudioStream

	out, err := cmd.Output()
	if err != nil {
		log.Println("Error:", err)
		return ffprobeOutput, ErrFileFormat
	}

	err = json.Unmarshal(out, &ffprobeOutput)
	if err != nil {
		log.Println("Error:", err)
		return ffprobeOutput, ErrFileFormat
	}

	if len(ffprobeOutput.Streams) == 0 {
		return ffprobeOutput, ErrFileFormat
	}

	audio := ffprobeOutput.Streams[0]

	SampleRate, err := strconv.Atoi(audio.SampleRate)
	if err != nil {
		log.Println("strconv SampleRate Error:", err)
		return ffprobeOutput, err
	}

	ffprobeOutput.Streams[0].sampleRate = SampleRate
	return ffprobeOutput, nil
}

func (b local) generateAudioEncodeMissoion(videofile string, AudioStream AudioStream) ([]worker.AudioFormat, error) {
	missions := []worker.AudioFormat{}
	if len(AudioStream.Streams) == 0 {
		return nil, ErrAudioFormat
	}

	audioStream := AudioStream.Streams[0]
	for i := range audioFormats {
		if audioStream.sampleRate < audioFormats[i].Hz {
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

func (b local) videoPath(videoID string) string {
	return fmt.Sprintf("%s%s.mp4", b.path, videoID)
}

func (b local) countFilesWithPrefix(videoID string) (int, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s%s_*.mp4", b.path, videoID))
	if err != nil {
		return 0, err
	}
	return len(files), nil
}
