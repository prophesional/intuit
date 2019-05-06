package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/prophesional/intuit"
)

func main() {
	h := Handler{}
	lambda.Start(h.Handle)
}

type Handler struct {
}

func (h *Handler) Handle(ctx context.Context, s3Event *events.S3Event) {
	path := "/tmp/tempfile"
	for _, r := range s3Event.Records {

		fmt.Println(r.S3.Bucket.Name)
		fmt.Println(r.S3.Object.Key)

		sess := session.Must(session.NewSession())
		fmt.Println("Session Created")
		// Create a downloader with the session and default options
		downloader := s3manager.NewDownloader(sess)
		// Create a file to write the S3 Object contents to.
		f, err := os.Create(path)

		if err != nil {
			fmt.Println("failed to create file")
		}
		// Write the contents of S3 Object to the file
		_, err = downloader.Download(f, &s3.GetObjectInput{
			Bucket: aws.String(r.S3.Bucket.Name),
			Key:    aws.String(r.S3.Object.Key),
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// print the content as 'bytes'

		players, err := intuit.ConvertToPlayer(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(players)

	}

}
