package utils

import (
	"app/backend/common"
	"path/filepath"

	"github.com/donnie4w/go-logger/logger"
)

var (
	Log = initLogger()
)

func initLogger() *logger.Logging {
	l := logger.NewLogger()
	l.SetOption(&logger.Option{

		Level:     logger.LEVEL_INFO,
		Console:   true,
		Format:    logger.FORMAT_LEVELFLAG | logger.FORMAT_SHORTFILENAME | logger.FORMAT_DATE | logger.FORMAT_MICROSECONDS,
		Formatter: "[{time}] {level} {file}: {message}\n",
		FileOption: &logger.FileTimeMode{ // 这里用时间切割
			Filename:   filepath.Join(common.WorkDir, "logs", common.AppName+".log"), // 日志文件路径
			Timemode:   logger.MODE_DAY,                                              // 按天
			Maxbuckup:  7,                                                            // 最多备份日志文件数
			IsCompress: false,                                                        // 是否压缩
		},
	})

	return l
}
