package global

import (
	"fmt"

	"github.com/spf13/viper"
)

// ConfigStruct 配置结构体，与config.yaml一一对应
type ConfigStruct struct {
	Server   ServerConfig   `mapstructure:"server"`
	Log      LogConfig      `mapstructure:"log"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type DatabaseConfig struct {
	Type            string `mapstructure:"type"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	Charset         string `mapstructure:"charset"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// InitConfig 初始化Viper配置
func InitConfig() error {
	v := viper.New()
	v.SetConfigFile("config/config.yaml") // 配置文件路径
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 绑定到全局Config变量
	if err := v.Unmarshal(&Config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}
