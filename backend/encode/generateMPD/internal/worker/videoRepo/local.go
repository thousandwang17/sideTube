/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 16:29:37
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-14 15:42:52
 * @FilePath: /encode/encodeVideo/handler.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package videoRepo

import (
	"context"
	"errors"
	"fmt"
	"generateMPD/internal/worker"
	"log"
	"os"
	"os/exec"
	"strings"
)

type local struct {
	path string
}

var (
	ErrDetct         = errors.New("ffprobe error")
	ErrFileNotFound  = errors.New("file not found")
	ErrFileFormat    = errors.New("ffempg failed tp get video format")
	ErrFormatMissing = errors.New("video or audio format missing")
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

func (b local) GenerateMPD(c context.Context, info worker.Mission, videoList, audioList []string) (mpdFileName, pngfileName, duration string, err error) {
	log.Println("generateMPDfile")

	mpdFileName, err = b.generateMPDfile(info.VideoId, videoList, audioList)
	if err != nil {
		return "", "", "", err
	}
	log.Println("createPng")

	pngfileName, err = b.createPng(info.VideoId)
	if err != nil {
		return "", "", "", err
	}
	log.Println("getVideoDuration")

	duration, err = b.getVideoDuration(info.VideoId)
	if err != nil {
		return "", "", "", err
	}
	log.Println("removeChunkVideo")

	b.removeChunkVideo(c, info)
	log.Println("down")

	return mpdFileName, pngfileName, duration, nil
}

func (b local) generateMPDfile(videoID string, videoList, audioList []string) (fileName string, err error) {

	// generate output file path
	outputFile := fmt.Sprintf("%s%s.mpd",
		b.path,
		videoID)

	// check  each list lenght must greater then 0
	if len(videoList) == 0 || len(audioList) == 0 {
		return "", ErrFormatMissing
	}

	// init cmd args
	index, input, copyMap, adaptation_sets := 0, []string{}, []string{}, strings.Builder{}

	// video list convet to cmd args
	adaptation_sets.WriteString("id=0,streams=")
	for i := range videoList {
		adaptation_sets.WriteString(fmt.Sprint(index))
		input = append(input, "-f", "webm_dash_manifest", "-i", videoList[i])
		copyMap = append(copyMap, "-map", fmt.Sprint(index))
		index++

		if i < len(videoList)-1 {
			adaptation_sets.WriteString(",")
		}
	}

	// audio list convet to cmd args
	adaptation_sets.WriteString(" id=1,streams=")
	for i := range audioList {
		adaptation_sets.WriteString(fmt.Sprint(index))
		input = append(input, "-f", "webm_dash_manifest", "-i", audioList[i])
		copyMap = append(copyMap, "-map", fmt.Sprint(index))
		index++

		if i < len(audioList)-1 {
			adaptation_sets.WriteString(",")
		}
	}

	// combine all args
	arg := []string{}
	arg = append(arg, input...)
	arg = append(arg, "-c", "copy")
	arg = append(arg, copyMap...)
	arg = append(arg, "-f", "webm_dash_manifest")
	arg = append(arg, "-y")
	arg = append(arg, "-adaptation_sets", adaptation_sets.String(), outputFile)

	cmd := exec.Command("ffmpeg", arg...)

	err = cmd.Run()
	if err != nil {
		log.Println("Error generate :", err, cmd.String())
		return "", ErrFileFormat
	}

	return fmt.Sprintf("%s.mpd", videoID), nil
}

func (b local) createPng(videoID string) (fileName string, err error) {

	// generate output file path
	inputFile := fmt.Sprintf("%s%s.mp4",
		b.path,
		videoID)

	// generate output file path
	outputFile := fmt.Sprintf("%s%s_default.png",
		b.path,
		videoID)

	startTime := "00:00"

	cmd := exec.Command("ffmpeg",
		"-ss", startTime,
		"-i", inputFile,
		"-vframes", "1",
		"-filter:v", "scale=854:-2",
		"-q:v", "1",
		"-y",
		outputFile,
	)

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error Png:", err)
		return "", err
	}

	return fmt.Sprintf("%s_default.png", videoID), nil
}

func (b local) getVideoDuration(videoID string) (string, error) {

	filePath := fmt.Sprintf("%s%s.mp4",
		b.path,
		videoID)

	// command to get the duration of the video
	cmd := exec.Command("ffprobe", "-i", filePath, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")

	// run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error Duration:", err)
		return "", err
	}

	// parse the duration from the output
	duration := strings.TrimSpace(string(output))

	// print the duration
	return duration, nil
}

func (b local) removeChunkVideo(c context.Context, info worker.Mission) {
	for part_id := 0; part_id < info.TotalChunk; part_id++ {
		if err := os.Remove(fmt.Sprintf("%s%s_%d.mp4",
			b.path,
			info.VideoId,
			part_id,
		)); err != nil {
			log.Printf("%s part: %d , copy err %v ", info.VideoId, part_id, err)
		}
	}
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
