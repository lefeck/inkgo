package service

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"inkgo/config"
	"inkgo/global"
)

type LoggerService struct {
	*config.LoggerConfig
}

func NewLoggerService(conf *config.LoggerConfig) *LoggerService {
	return &LoggerService{conf}
}

func (loggerService *LoggerService) WriteLog() error {
	writeSyncer := getLogWriter(loggerService.Filename, loggerService.MaxSize, loggerService.MaxBackups, loggerService.MaxAge)
	encoder := getEncoder()
	var log = new(zapcore.Level)
	err := log.UnmarshalText([]byte(loggerService.Level))
	if err != nil {
		return err
	}

	core := zapcore.NewCore(encoder, writeSyncer, log)

	global.Log = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(global.Log) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
