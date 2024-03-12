package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 设置静态资源和HTML模板文件夹的相对路径
	r.Static("/static", "static")
	r.Static("/assets", "assets")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/videoview", func(c *gin.Context) {
		c.HTML(200, "videoview.html", nil)
	})

	// 处理视频上传
	r.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("video")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请选择一个视频文件"})
			return
		}

		title := c.PostForm("title")
		if title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请输入视频标题"})
			return
		}

		// 检查文件类型
		if header.Header.Get("Content-Type") != "video/mp4" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持MP4格式的视频文件"})
			return
		}

		// 检查文件大小
		if header.Size > 200*1024*1024 { // 假设最大文件大小为100MB
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过100MB"})
			return
		}

		// 生成新的文件名
		ext := filepath.Ext(header.Filename)
		newFileName := fmt.Sprintf("video-%d-%s%s", time.Now().UnixNano(), title, ext)
		filePath := filepath.Join("static/videos", newFileName)

		out, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "视频上传成功",
			"title":   title,
		})
	})

	// 处理视频列表请求
	r.GET("/videos", func(c *gin.Context) {
		videoDirectory := "static/videos"
		files, err := os.ReadDir(videoDirectory)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var videoFiles []map[string]string

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			ext := filepath.Ext(file.Name())
			if ext == ".mp4" {
				videoFiles = append(videoFiles, map[string]string{
					"title": filepath.Base(file.Name()[:len(file.Name())-len(ext)]),
					"url":   "/static/videos/" + file.Name(),
				})
			}
		}

		c.JSON(http.StatusOK, videoFiles)
	})

	// 添加路由以处理单独播放视频的页面
	r.GET("/videoplay/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.HTML(http.StatusOK, "videoplay.html", gin.H{"filename": filename})
	})

	// 启动服务器
	r.Run(":8080")
}
