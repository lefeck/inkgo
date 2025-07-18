package service

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"inkgo/config"
	"mime/multipart"
)

// 文件上传
type UploadService struct {
	Conf *config.StorageConfig
}

type IUploadService interface {
	UploadFile(file multipart.File, fileSize int64) (string, error)
}

func NewUploadService(conf *config.StorageConfig) IUploadService {
	return &UploadService{conf}
}

func (us *UploadService) UploadFile(file multipart.File, fileSize int64) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: us.Conf.Bucket,
	}
	mac := qbox.NewMac(us.Conf.AccessKey, us.Conf.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		return "", err
	}
	url := us.Conf.StorageUrl + "/" + ret.Key
	return url, nil
}
