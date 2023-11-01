package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

//	func main() {
//		InitLogger()
//		defer logger.Sync() //将缓冲区的内容刷到磁盘上
//		for i := 0; i < 40000; i++ {
//			logger.Info("test for log rotate")
//		}
//
//		//simpleHttpGet("www.baidu.com")
//		simpleHttpGet("http://www.baidu.com")
//	}
func main() {
	//r := gin.Default()
	InitLogger()
	r := gin.New()
	r.Use(GinLogger(logger), GinRecovery(logger, true))
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello liwenzhou.com!")
	})
	r.Run(":8080")
}

func InitLogger() {
	//自定义日志配置信息
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	//创建自定义日志
	logger = zap.New(core, zap.AddCaller()) //zap.AddCaller()  添加调用函数信息
	sugarLogger = logger.Sugar()
}

// 编码格式-以json格式记录日志
func getEncoder() zapcore.Encoder {
	encodeconfig := zap.NewProductionEncoderConfig()
	encodeconfig.EncodeTime = zapcore.ISO8601TimeEncoder //易读性更好的时间
	return zapcore.NewJSONEncoder(encodeconfig)
	//return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
}

//func getLogWriter() zapcore.WriteSyncer {
//	//创建文件方便写入日志
//	//file, _ := os.Create("./test.log")
//	//直接读写文件
//	file, _ := os.OpenFile("./test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
//	return zapcore.AddSync(file)
//}

// 第三方库实现日志切割
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,     // MB
		MaxBackups: 5,     //	最大备份数量
		MaxAge:     30,    //最大保存时间
		Compress:   false, //是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		logger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
