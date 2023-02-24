/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 15:46:39
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-24 17:09:24
 * @FilePath: /jwtGenerate/internal/jwtGenerate/video.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package jwtGenerate

const (
	timeFormart = "2006-01-02 15:04:05"
)

// RefreshToken have long ttl time (7d , 30d ...).
// RefreshToken used for request new ApiToken.
// jwt Format
type AccessToken string

// ApiToken have short ttl time (15 mins, 30min...).
// ApiToken is used for request server.
// jwt Format
type RefreshToken string

// Login need to return both kind of token
type LoginToken struct {
	AT AccessToken
	RT RefreshToken
}

type UserInfo struct {
	UserId   string `json:"user_id" bson:"userId"`
	UserName string `json:"user_name" bson:"userName"`
}
