/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-30 13:24:39
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-06 14:50:17
 * @FilePath: /videoUpload/internal/common/awsS3/session.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package awsS3

import "github.com/aws/aws-sdk-go/aws/session"

var awsSession *session.Session

func init() {
	awsSession = session.Must(session.NewSession())
}

func GetAWSSession() *session.Session {
	return awsSession
}
