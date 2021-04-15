package server

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
	"net/http"
)

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func PutFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	fmt.Println(r.Header.Get("Content-Type"))
	file, handler, err := r.FormFile("test")
	if err != nil {
		fmt.Println("Error retrieving File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Configuration error, " + err.Error())
	}
	client := s3.NewFromConfig(cfg)
	fmt.Println(cfg.Region)
	bucket := viper.GetString("S3bucket")
	filename := handler.Filename
	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &filename,
		Body:   file,
	}
	_, err = PutFile(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Uploading error")
		fmt.Println(err.Error())
		return
	}

}
