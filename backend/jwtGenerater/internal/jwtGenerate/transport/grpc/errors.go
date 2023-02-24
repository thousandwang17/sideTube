/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-24 21:54:42
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-24 21:55:01
 * @FilePath: /jwtGenerater/internal/jwtGenerate/transport/grpc/errors.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package grpc

import "errors"

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
