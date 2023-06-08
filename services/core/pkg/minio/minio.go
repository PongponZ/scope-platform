package minio

import (
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinio(endpoint string, accessKey string, secert string) *minio.Client {
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secert, ""),
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("minio connected :: %s", endpoint)
	return client
}
