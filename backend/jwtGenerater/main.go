/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-30 13:22:07
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-24 14:18:22
 * @FilePath: /jwtGenerate/cmd/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	_ "sideTube/jwtGenerate/internal/common/mongodb"
)

func main() {
	startHttpServer()
}
