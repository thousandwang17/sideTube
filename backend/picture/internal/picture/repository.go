/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-02 16:15:59
 * @FilePath: /picture/internal/picture/repository.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package picture

import (
	"context"
	"io"
)

type PictureRepository interface {
	GetPicture(c context.Context, videoId string) (data io.ReadCloser, err error)
}
