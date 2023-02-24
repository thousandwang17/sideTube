package videoRepository

import (
	"context"
	"errors"
	"io"
	"os"
	"sideTube/videoUpload/internal/videoUpload"
	v "sideTube/videoUpload/internal/videoUpload"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
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
func NewS3(AWSSession *session.Session) v.VideoRepository {
	svc := s3.New(AWSSession, aws.NewConfig().WithRegion("ap-northeast-3"))
	return bucket{
		svc,
		os.Getenv("aws_video_bucket"),
	}
}

/**
 * @description:
 * @param {context.Context} c
 * @return {*}
 */
func (b bucket) CreateMultipartUpload(c context.Context, id string) (v.VideoRepoMeta, error) {

	if exist, kerr := b.keyExists(id); kerr == nil {
		return v.VideoRepoMeta{}, kerr
	} else if exist {
		return v.VideoRepoMeta{}, errors.New("videoId is already exists")
	}

	// using aws s3 sdk
	respond, err := b.svc.CreateMultipartUploadWithContext(c, &s3.CreateMultipartUploadInput{
		Bucket: &b.name,
		Key:    &id,
	})

	if err != nil {
		return v.VideoRepoMeta{}, err
	}

	return v.VideoRepoMeta{
			Id:       id,
			UploadId: *respond.UploadId,
		},
		nil
}

/**
 * @description:
 * @param {context.Context} c
 * @param {string} id
 * @param {int64} partId
 * @param {io.ReadSeeker} file
 * @return {*}
 */
func (b bucket) UploadPart(c context.Context, v videoUpload.VideoRepoMeta, file io.ReadSeeker) error {

	// using aws s3 sdk
	_, err := b.svc.UploadPartWithContext(c, &s3.UploadPartInput{
		Bucket:     &b.name,
		Key:        aws.String(v.Id),
		Body:       file,
		PartNumber: aws.Int64(v.PartID),
		UploadId:   aws.String(v.UploadId),
	})

	if err != nil {
		return err
	}

	return nil
}

/**
 * @description:
 * @param {context.Context} c
 * @param {string} id
 * @return {*}
 */
func (b bucket) CompleteMultipartUpload(c context.Context, id string) error {

	_, err := b.svc.CompleteMultipartUploadWithContext(c, &s3.CompleteMultipartUploadInput{
		Bucket: &b.name,
		Key:    &id,
	})

	if err == nil {
		return nil
	}

	// may be already upload
	if err.Error() == s3.ErrCodeNoSuchUpload {
		exists, berr := b.keyExists(id)

		if exists && berr == nil {
			return nil
		} else if berr != nil {
			return berr
		}

		return errors.New("video not found")
	}

	return err
}

/**
 * @description:
 * @param {context.Context} c
 * @param {string} id
 * @return {*}
 */
func (b bucket) AbortUpload(c context.Context, id string) error {

	// using aws s3 sdk
	_, err := b.svc.AbortMultipartUploadWithContext(c, &s3.AbortMultipartUploadInput{
		Bucket: &b.name,
		Key:    &id,
	})

	if err != nil {
		return err
	}

	return nil
}

/**
 * @description:
 * @param {string} key
 * @return {*}
 */
func (b bucket) keyExists(key string) (bool, error) {

	_, err := b.svc.HeadObject(&s3.HeadObjectInput{
		Bucket: &b.name,
		Key:    &key,
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}

	return true, nil
}
