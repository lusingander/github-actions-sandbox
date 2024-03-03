package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

func run() error {

	// setup localstack image
	ctx := context.Background()
	container, err := localstack.RunContainer(
		ctx,
		testcontainers.WithImage("localstack/localstack:3.2.0"),
		testcontainers.CustomizeRequest(
			testcontainers.GenericContainerRequest{
				ContainerRequest: testcontainers.ContainerRequest{
					Env: map[string]string{
						"SERVICES": "s3",
					},
				},
			},
		),
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// get access parameters
	port, err := container.MappedPort(ctx, nat.Port("4566/tcp"))
	if err != nil {
		return err
	}

	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		return err
	}
	defer provider.Close()

	host, err := provider.DaemonHost(ctx)
	if err != nil {
		return err
	}

	// setup s3 client
	url := fmt.Sprintf("http://%s:%d", host, port.Int())
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, opts ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           url,
				SigningRegion: region,
			}, nil
		})
	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy"),
		),
	)
	if err != nil {
		return err
	}
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	// setup fixtures
	buckets := []string{
		"test-bucket-1",
		"test-bucket-2",
		"test-bucket-3",
	}
	for _, bucket := range buckets {
		_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return err
		}
		log.Printf("Create bucket success: name = %v\n", bucket)
	}

	// wait :)
	// log.Printf("Start %v: Wait...\n", url)
	// fmt.Scanln()

	// check
	output, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}
	log.Printf("[Created buckets]\n")
	for _, bucket := range output.Buckets {
		log.Printf("name = %v, creation_date = %v\n", *bucket.Name, bucket.CreationDate)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
