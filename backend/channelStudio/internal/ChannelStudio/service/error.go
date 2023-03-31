/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-09 15:55:34
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 18:05:32
 * @FilePath: /channelStudio/internal/ChannelStudio/service/error.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotReachMinWidth = errors.New(fmt.Sprintf("not reach the min width require : %d", minPicturWidth))
	ErrHigthWidthRatio  = errors.New(fmt.Sprintf("image h/w ratio should between %f ~  %f", maxPicturRatio, minPicturRatio))
	ErrVaild            = errors.New("args vaild failed")
	ErrSetPermission    = errors.New("title and desc are empty")
)

type serviceErr struct {
	E error
}

func (m serviceErr) Error() string {
	return m.E.Error()
}

func (m serviceErr) StatusCode() int {
	return http.StatusBadRequest
}
