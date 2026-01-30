/*
 * Copyright 2025 Bronya0 <tangssst@163.com>.
 * Author Github: https://github.com/Bronya0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	// Version 会在编译时注入 -ldflags="-X 'app/backend/common.Version=${{ github.event.release.tag_name }}'"
	Version = ""
)

const (
	AppName     = "API-King"
	Desc        = ""
	Width       = 600
	Height      = 500
	Theme       = "dark"
	LocalDBFile = "API-King.db"
	LogPath     = "API-King.log"
	Language    = "zh-CN"
)

const ()

var (
	WorkDir    = getAppDataDir()
	HttpClient = initHttpClient()
)

// 获取并创建应用数据目录,mac上强制限制
func getAppDataDir() string {
	// 获取用户配置目录
	// macOS (Sandbox) -> ~/Library/Containers/com.xx.xx/Data/Library/Application Support
	// macOS (Dev)     -> ~/Library/Application Support
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Error getting user config dir: %v", err)
		return ""
	}

	// 拼接 App 名称
	appDir := filepath.Join(configDir, AppName)

	// 确保目录存在，不存在则创建
	if err := os.MkdirAll(appDir, 0755); err != nil {
		log.Printf("Error creating app data dir: %v", err)
		return ""
	}

	return appDir
}

func initHttpClient() *resty.Client {
	// 创建一个新的resty客户端实例
	return resty.New().
		// 设置HTTP传输配置
		SetTransport(&http.Transport{
			MaxIdleConns:        10,               // 最大空闲连接数
			MaxIdleConnsPerHost: 10,               // 每个主机的最大空闲连接数
			MaxConnsPerHost:     20,               // 每个主机的最大连接数
			IdleConnTimeout:     60 * time.Second, // 空闲连接超时
			DisableKeepAlives:   false,            // 启用 Keep-Alive
		}).
		// 设置客户端超时时间
		SetTimeout(60 * time.Second).
		// 设置TLS客户端配置，跳过TLS验证
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}
