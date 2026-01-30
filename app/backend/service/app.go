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

package service

import (
	"app/backend/common"
	"app/backend/model"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/wailsapp/wails/v3/pkg/application"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// Setting 设置表模型
type Setting struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

// DBStats 数据库统计信息
type DBStats struct {
	TotalRows int64  `json:"total_rows"`
	FileSize  int64  `json:"file_size"`
	SizeStr   string `json:"size_str"`
}

// 配置默认值常量，统一管理，避免前后端魔法数字
const (
	DefaultTheme = "dark"
)

// App 应用程序结构体
type App struct {
	app    *application.App
	db     *gorm.DB
	dbPath string
	window *application.WebviewWindow
}

// NewApp 创建应用实例
func NewApp() *App {
	dbPath := filepath.Join(common.WorkDir, common.AppName+".db")
	dsn := dbPath + "?cache=shared&mode=rwc&_journal_mode=WAL"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // Info 会打印 SQL
				Colorful:      true,
			},
		),
	})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(0)
	sqlDB.SetConnMaxLifetime(0)

	app := &App{
		db:     db,
		dbPath: dbPath,
	}

	return app
}

// SetApp 设置应用实例
func (a *App) SetApp(app *application.App) {
	a.app = app
}

// SetWindow 更新窗口实例
func (a *App) SetWindow(window *application.WebviewWindow) {
	a.window = window
}

// ServiceStartup 应用启动时执行 (v3 Service Lifecycle)
func (a *App) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {

	// InitializeServices 初始化核心后端服务
	a.InitializeServices()

	return nil
}

// InitializeServices 初始化核心后端服务
func (a *App) InitializeServices() {

	log.Println("=== 后端服务初始化 ===")

	a.Migrate()

	// 首次启动初始化默认配置
	a.initDefaultSettings()

	// 启动定时自动优化数据库
	go a.startAutoVacuum()

}

func (a *App) ServiceShutdown() error {
	return nil
}

// startAutoVacuum 开启定时自动优化数据库任务
func (a *App) startAutoVacuum() {
	// 1. 启动后延迟 1 分钟执行第一次优化，避免与启动时的 IO 抢占资源
	time.AfterFunc(1*time.Minute, func() {
		a.runVacuum()
	})

	// 2. 每隔 12 小时执行一次深度优化
	ticker := time.NewTicker(12 * time.Hour)
	for range ticker.C {
		a.runVacuum()
	}
}

// runVacuum 执行数据库压缩优化
func (a *App) runVacuum() {
	log.Println("开始后台执行数据库 VACUUM 优化...")
	start := time.Now()
	// VACUUM 会重建整个数据库文件，释放未使用的空间并整理碎片
	if err := a.db.Exec("VACUUM").Error; err != nil {
		log.Printf("后台 VACUUM 失败: %v", err)
	} else {
		log.Printf("后台 VACUUM 完成，耗时: %v", time.Since(start))
	}
}

// ShowWindow 显示应用窗口
func (a *App) ShowWindow() {
	if a.window != nil {
		a.window.Show()
	}
}

// HideWindow 隐藏应用窗口
func (a *App) HideWindow() {
	log.Println("HideWindow")
	a.window.Hide()
}

// SetWindowAlwaysOnTop 设置窗口是否置顶
func (a *App) SetWindowAlwaysOnTop(alwaysOnTop bool) {
	if a.window != nil {
		a.window.SetAlwaysOnTop(alwaysOnTop)
	}
}

// QuitApp 退出应用
func (a *App) QuitApp() {
	a.app.Quit()
}

// SetTheme 供 JS 调用
// theme 参数可取值: "dark", "light", "system"
func (a *App) SetTheme(theme string) {
	application.InvokeAsync(func() {
		switch theme {
		case "dark":
			a.window.SetBackgroundColour(application.NewRGBA(30, 30, 30, 255))
		case "light":
			a.window.SetBackgroundColour(application.NewRGBA(255, 255, 255, 255))
		case "system":
		}
	})
}

// GetDBDir 获取数据库存储绝对路径（目录）
func (a *App) GetDBDir() string {
	return filepath.Dir(a.dbPath)
}

// OpenDirectory 打开指定目录
func (a *App) OpenDirectory(path string) error {
	// 使用 Wails v3 内置 API 打开文件管理器
	// 第二个参数为 true 表示在文件管理器中选中该文件/目录（如果支持）
	// 这里我们只需要打开目录，所以传 false
	err := a.app.Env.OpenFileManager(path, false)
	if err != nil {
		log.Printf("无法打开目录 %s: %v", path, err)
		return err
	}
	return nil
}

// GetConfig 获取配置
func (a *App) GetConfig() map[string]interface{} {
	// 主题没有配置时使用统一的默认值
	theme := a.GetSetting("theme")
	if theme == "" {
		theme = DefaultTheme
	}

	return map[string]interface{}{
		"theme":     theme,
		"autoStart": a.GetSetting("auto_start") == "true",
	}
}

// GetVersion returns the application version
func (a *App) GetVersion() string {
	return common.Version
}

func (a *App) GetAppName() string {
	return common.AppName
}
func (a *App) GetDesc() string {
	return common.Desc
}

// getIntSetting 获取整数配置，带默认值
func (a *App) getIntSetting(key string, defaultValue int) int {
	value := a.GetSetting(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

// GetSetting 从数据库获取设置
func (a *App) GetSetting(key string) string {
	var s Setting
	if err := a.db.Where("key = ?", key).First(&s).Error; err != nil {
		return ""
	}
	return s.Value
}

// SaveSetting 保存设置到数据库
func (a *App) SaveSetting(key string, value string) {
	a.db.Save(&Setting{Key: key, Value: value})
}

func formatSize(size int64) string {
	if size < 0 {
		return "-" + formatSize(-size)
	}
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// Migrate 数据库迁移
func (a *App) Migrate() {
	a.db.Exec("PRAGMA foreign_keys = ON;")

	// 1. 使用 GORM 自动迁移标准表结构（处理建表、新增列、索引）
	if err := a.db.AutoMigrate(&Setting{}, &model.History{}, &model.ApiRequest{}); err != nil {
		log.Print("数据库 AutoMigrate 失败:", err)
	}
}

// initDefaultSettings 初始化默认配置
func (a *App) initDefaultSettings() {
	// 定义需要初始化的默认配置项
	defaults := map[string]string{
		"theme":      DefaultTheme,
		"auto_start": "false",
	}

	settings := make([]Setting, 0, len(defaults))
	for k, v := range defaults {
		settings = append(settings, Setting{Key: k, Value: v})
	}
	// 批量插入，重复 key 忽略
	a.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&settings)
}
