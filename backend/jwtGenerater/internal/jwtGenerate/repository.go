/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-24 17:09:59
 * @FilePath: /jwtGenerate/internal/jwtGenerate/repository.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package jwtGenerate

import (
	"context"
)

type MetaRepository interface {
	LogInCheck(c context.Context, userId, passWordSH256 string) (UserInfo, error)
}
