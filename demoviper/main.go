package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// 序列化需要加mapstructure
type Config struct {
	Port        int    `mapstructure:"host"`
	Version     string `mapstructure:"version"`
	MysqlConfig `mapstructure:"mysql"`
}

type MysqlConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Dbname string `mapstructure:"dbname"`
}

type QQQ struct {
	Version string `mapstructure:"version"`
	QAQ     `mapstructure:"qaq"`
}
type QAQ struct {
	QWER string `mapstructure:"qwer"`
	ASDF string `mapstructure:"asdf"`
}

func main() {
	//读入配置文件
	viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")        // 配置文件名称(无扩展名)
	viper.SetConfigType("json")          // 如果配置文件的名称中没有扩展名，则需要配置此项，用于从远端服务器如etcd读取指定类型的文件，本地不行
	//viper.AddConfigPath("/etc/appname/")  // 查找配置文件所在的路径
	//viper.AddConfigPath("$HOME/.appname") // 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")    // 还可以在工作目录中查找配置
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	//实时监控配置文件的改变
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})

	var c QQQ
	err = viper.Unmarshal(&c)
	if err != nil {
		fmt.Printf("unmarshall error : %v", err)
	}

	fmt.Printf("c : %#v", c)

	//r := gin.Default()
	//r.GET("/hello", func(c *gin.Context) {
	//	c.String(http.StatusOK, viper.GetString("version"))
	//})
	//r.Run()
}
