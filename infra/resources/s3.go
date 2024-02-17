package resource

import (
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewBucket(name string) s3.Bucket {
	return s3.NewBucket(r.S, jsii.String(name), &s3.BucketProps{
		BucketName: jsii.String(name),
	})
}

func (r *ResourceService) GetBucketFromName(name string) s3.IBucket {
	return s3.Bucket_FromBucketName(r.S, jsii.String(name), jsii.String(name))
}
