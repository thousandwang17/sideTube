/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-04 13:42:46
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-02 16:59:10
 * @FilePath: /recommend/internal/common/simpleKit/endpoint.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package endpoint

import (
	"context"
)

type EndPoint func(context.Context, interface{}) (interface{}, error)

type MiddleWare func(EndPoint) EndPoint
