/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-15 16:29:37
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-28 16:08:09
 * @FilePath: /encode/encodeVideo/handler.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package videoRepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mergeUploadFile/internal/worker"
	"os"
	"os/exec"
	"strings"
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
func NewLoacl(path string) worker.VideoRepository {
	return local{
		path: path,
	}
}

/**
 * @description:
 * @param {context.Context} c
 * @param {string} id
 * @return {*}
 */
func (b local) MergeEncodeVideo(c context.Context, info worker.Mission) (string, error) {
	mergeFileName := fmt.Sprintf("%s%s.%d.%v.tmp.webm",
		b.path,
		info.VideoId,
		info.Height,
		info.Fps,
	)

	// Create the input file list as a string and write to a temporary file
	var inputStr strings.Builder
	for i := 0; i < info.TotalChunk; i++ {
		fmt.Fprintf(&inputStr, "file '%s%s_%d.%d.%v.webm'\n",
			b.path,
			info.VideoId,
			i,
			info.Height,
			info.Fps)
	}

	tmpfile, err := os.CreateTemp("", "input*.txt")
	if err != nil {
		log.Println("Faled to Temp file ", tmpfile)
		return "", err
	}

	defer os.Remove(tmpfile.Name())
	_, err = tmpfile.WriteString(inputStr.String())
	if err != nil {
		log.Println("Faled to write to tmp file ", tmpfile)
		return "", err
	}

	err = tmpfile.Close()
	if err != nil {
		log.Println("Faled to close to tmp file ", tmpfile)
	}

	// Create the command
	cmd := exec.Command("ffmpeg",
		"-y",
		"-f", "concat",
		"-safe", "0",
		"-i", tmpfile.Name(),
		"-c", "copy", mergeFileName)
	log.Println(cmd.String())
	// Run the command and check for errors
	err = cmd.Run()
	if err != nil {
		log.Println("Faled to merge to video file ", info.VideoId, info.MissionID, err)
		return "", err
	}

	dashFileName := fmt.Sprintf("%s%s.%d.%v.webm",
		b.path,
		info.VideoId,
		info.Height,
		info.Fps,
	)

	// re-dash the merged file with copy method
	cmd = exec.Command("ffmpeg",
		"-i", mergeFileName,
		"-c:v", "copy",
		"-an",
		"-f", "webm",
		"-dash", "1",
		dashFileName,
	)
	log.Println(cmd.String())
	err = cmd.Run()
	if err != nil {
		log.Println("Faled to dash the merged video", info.VideoId, info.MissionID, err)
		return "", err
	}

	b.removeChunkFile(c, info)
	os.Remove(mergeFileName)

	log.Printf("%s.%d.%v.dash.webm  dash done!!", info.VideoId, info.Height, info.Fps)
	return dashFileName, nil
}

func (b local) removeChunkFile(c context.Context, info worker.Mission) {
	for part_id := 0; part_id < info.TotalChunk; part_id++ {
		if err := os.Remove(fmt.Sprintf("%s%s_%d.%d.%v.webm",
			b.path,
			info.VideoId,
			part_id,
			info.Height,
			info.Fps)); err != nil {
			log.Printf("%s part: %d , copy err %v ", info.VideoId, part_id, err)
		}
	}
}

func (b local) checkFileStat(fileName string) (bool, error) {
	if _, err := os.Stat(fmt.Sprintf("%s%s", b.path, fileName)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
