package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

func Init() (err error) {
	//相对路径，相对于可执行文件的路径
	//viper.SetConfigFile("./config.yaml")

	////通过输入参数设置config文件路径
	//viper.SetConfigFile(filename)

	//绝对路径,当前系统的文件路径
	//viper.SetConfigFile("D:\\GOOOOOOOOOOOOOOOOOOOOOOOO\\go_web\\config.yaml")
	//读入配置文件
	viper.SetConfigName("config") // 配置文件名称(无扩展名)

	//从远程配置如etcd中心的获取的文件类型
	viper.SetConfigType("yaml") // 如果配置文件的名称中没有扩展名，则需要配置此项

	viper.AddConfigPath(".")   // 设置从当前路径找（使用相对路径）
	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {            // 处理读取配置文件的错误
		fmt.Printf("get config error : %v", err)
		return
	}

	if err = viper.Unmarshal(&Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})

	//使用结构体保存，每次改配置后，要重新反序列化
	//viper.WatchConfig()
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	fmt.Println("配置文件修改了...")
	//	if err := viper.Unmarshal(Conf); err != nil {
	//		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	//	}
	//})
	return
}
