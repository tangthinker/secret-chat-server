package oss

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/tangthinker/secret-chat-server/core"
)

type Service struct {
	storagePath string
	accessUrl   string
	cleanTtl    time.Duration
}

func NewService() *Service {
	storagePath := core.GlobalHelper.Config.GetString("oss.storage-path")
	accessUrl := core.GlobalHelper.Config.GetString("oss.access-url")
	cleanTtl := core.GlobalHelper.Config.GetDuration("oss.clean-ttl")
	s := &Service{
		storagePath: storagePath,
		accessUrl:   accessUrl,
		cleanTtl:    cleanTtl,
	}
	s.startCleanTask()
	return s
}

// Upload 上传文件 返回文件的url
func (s *Service) Upload(ctx context.Context, reader io.Reader, fileName string) (string, error) {
	ext := filepath.Ext(fileName)
	// 文件名格式为：uuid_时间戳.ext
	filename := strings.ReplaceAll(strings.ToLower(uuid.New().String()), "-", "") + "_" + time.Now().Format("20060102150405") + ext
	filePath := filepath.Join(s.storagePath, filename)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("create directory error: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("create file error: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return "", fmt.Errorf("copy file error: %v", err)
	}

	accessUrl := s.accessUrl
	if !strings.HasSuffix(accessUrl, "/") {
		accessUrl += "/"
	}
	return accessUrl + filename, nil
}

func (s *Service) startCleanTask() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("startCleanTask error: %v", err)
			}
		}()
		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			s.clean()
		}
	}()
}

func (s *Service) clean() {
	files, err := os.ReadDir(s.storagePath)
	if err != nil {
		return
	}
	cleanedFiles := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		fileTime, err := extractFileTime(filename)
		if err != nil {
			continue
		}
		if time.Since(*fileTime) > s.cleanTtl {
			if err := os.Remove(filepath.Join(s.storagePath, filename)); err != nil {
				log.Errorf("clean: remove file error: %v", err)
			}
		}
		cleanedFiles = append(cleanedFiles, filename)
	}
	if len(cleanedFiles) > 0 {
		log.Infof("clean: cleaned %d files", len(cleanedFiles))
	}
}

// extractFileTime 提取文件时间
func extractFileTime(filename string) (*time.Time, error) {
	parts := strings.Split(filename, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid filename: %s", filename)
	}
	parts = strings.Split(parts[1], ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid filename: %s", filename)
	}
	timeStr := parts[0]
	time, err := time.Parse("20060102150405", timeStr)
	if err != nil {
		return nil, fmt.Errorf("parse time error: %v", err)
	}
	return &time, nil
}
