package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/mock"

	// has purpose as of method 1, in mock q&a section, readme
	//_ "github.com/golang/mock/mockgen/model"

	"github.com/jaytaylor/mockery-example/mocks"
)

func main() {
	mockS3 := &mocks.S3API{}
	//mockS3 := &mocks.MockS3API{}

	mockResultFn := func(input *s3.ListObjectsInput) *s3.ListObjectsOutput {
		output := &s3.ListObjectsOutput{}
		output.SetCommonPrefixes([]*s3.CommonPrefix{
			{
				Prefix: aws.String("2017-01-01"),
			},
		})
		output.SetContents([]*s3.Object{
			{
				Key: aws.String("foo-object"),
			},
		})
		return output
	}

	mockS3.On("ListObjects",
		mock.MatchedBy(func(input *s3.ListObjectsInput) bool {
			return input.Delimiter != nil && *input.Delimiter == "/" && input.Prefix == nil
		}),
	).
		Return(mockResultFn, nil)

	// action
	listingInput := &s3.ListObjectsInput{
		Bucket:    aws.String("foo"),
		Delimiter: aws.String("/"),
	}
	listingOutput, err := mockS3.ListObjects(listingInput)
	if err != nil {
		panic(err)
	}

	// print
	for _, x := range listingOutput.CommonPrefixes {
		fmt.Printf("common prefix: %+v\n", *x)
	}
	for _, x := range listingOutput.Contents {
		fmt.Printf("content: %+v\n", *x)
	}
}
