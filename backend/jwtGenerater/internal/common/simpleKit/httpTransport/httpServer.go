/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-04 13:45:08
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-03 14:49:51
 * @FilePath: /jwtGenerate/internal/common/simpleKit/httpTransport/httpServer.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package httptransport

import (
	"context"
	"log"
	"net/http"
	p "sideTube/jwtGenerate/internal/common/simpleKit/endpoint"
)

type DecoderFuc func(context.Context, *http.Request) (interface{}, error)
type EncoderFuc func(context.Context, http.ResponseWriter, interface{}) error

type ServerBefore func(context.Context, *http.Request) context.Context

type ServerOprion func(s *HttpTransport)

type HttpTransport struct {
	endpoint p.EndPoint
	decode   DecoderFuc
	encode   EncoderFuc
	before   []ServerBefore
}

func NewHttpTransport(e p.EndPoint, de DecoderFuc, en EncoderFuc, options ...ServerOprion) *HttpTransport {
	s := &HttpTransport{
		endpoint: e,
		decode:   de,
		encode:   en,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func NewServerBefore(b ...ServerBefore) ServerOprion {
	return func(s *HttpTransport) {
		s.before = append(s.before, b...)
	}
}

func (h HttpTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	for _, b := range h.before {
		ctx = b(ctx, r)
	}

	request, err := h.decode(ctx, r)
	if err != nil {
		errorHandle(ctx, err)
		errorEncode(ctx, err, w)
		return
	}

	data, err := h.endpoint(ctx, request)
	if err != nil {
		errorHandle(ctx, err)
		errorEncode(ctx, err, w)
		return
	}

	err = h.encode(ctx, w, data)
	if err != nil {
		errorHandle(ctx, err)
		errorEncode(ctx, err, w)
		return
	}
}

func errorHandle(ctx context.Context, err error) {
	log.Println("httpTransport error: ", err)
}

func errorEncode(_ context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("Content-Type", contentType)

	code := http.StatusInternalServerError
	if status, ok := err.(StatusCoder); ok {
		code = status.StatusCode()
	}
	w.WriteHeader(code)
	w.Write(body)
}

type StatusCoder interface {
	StatusCode() int
}
