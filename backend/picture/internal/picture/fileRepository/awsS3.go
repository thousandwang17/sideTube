/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-01 20:03:11
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-02 16:14:20
 * @FilePath: /picture/internal/picture/videoRepository/awsS3.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package fileRepository

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

type bucket struct {
	svc  *s3.S3
	name string
}

/**
 * @description:
 * @return {*}
 */
// func NewS3(AWSSession *session.Session) v.VideoRepository {
// 	svc := s3.New(AWSSession, aws.NewConfig().WithRegion("ap-northeast-3"))
// 	return bucket{
// 		svc,
// 		os.Getenv("aws_video_bucket"),
// 	}
// }
