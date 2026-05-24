package config

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Static StaticConfig `mapstructure:"static"`
	Redis  RedisConfig  `mapstructure:"redis"`
	Video  VideoConfig  `mapstructure:"video"`
	Expire ExpireConfig `mapstructure:"expire"`
}

type ServerConfig struct {
	Address string `mapstructure:"address"`
}

type MySQLConfig struct {
	DSN string `mapstructure:"dsn"`
}

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	RefreshExpire int64  `mapstructure:"refresh_expire"`
	MFAExpire     int64  `mapstructure:"mfa_expire"`
	MFAPeriod     int    `mapstructure:"mfa_period"`
}

type StaticConfig struct {
	BaseURL       string `mapstructure:"base_url"`
	UploadDir     string `mapstructure:"upload_dir"`
	AvatarPath    string `mapstructure:"avatar_path"`
	VideoPath     string `mapstructure:"video_path"`
	DefaultAvatar string `mapstructure:"default_avatar"`
}

type RedisConfig struct {
	Addr       string `mapstructure:"addr"`
	Password   string `mapstructure:"password"`
	DB         int    `mapstructure:"db"`
	ChatExpire int64  `mapstructure:"chat_expire"`
}

type VideoConfig struct {
	FeedLimit    int `mapstructure:"feed_limit"`
	PopularLimit int `mapstructure:"popular_limit"`
}

type ExpireConfig struct {
	RefreshToken int64 `mapstructure:"refresh_token"`
	MFA          int64 `mapstructure:"mfa"`
	Chat         int64 `mapstructure:"chat"`
}

var (
	Conf              *Config
	viperV            *viper.Viper
	rwMutex           sync.RWMutex
	onChangeCallbacks []func(*Config)  // 配置变更回调列表
)

func InitConfig() {
	viperV = viper.New()
	viperV.SetConfigName("config")
	viperV.SetConfigType("yaml")
	viperV.AddConfigPath("config")
	viperV.AddConfigPath("./config")
	viperV.AddConfigPath("../config")

	viperV.SetDefault("server.address", "0.0.0.0:8888")
	viperV.SetDefault("jwt.refresh_expire", 604800)
	viperV.SetDefault("jwt.mfa_expire", 604800)
	viperV.SetDefault("jwt.mfa_period", 30)
	viperV.SetDefault("static.base_url", "http://127.0.0.1:8888")
	viperV.SetDefault("static.avatar_path", "/static/avatars/")
	viperV.SetDefault("static.video_path", "/static/videos/")
	viperV.SetDefault("static.default_avatar", "/static/avatars/default_avatar.png")
	viperV.SetDefault("redis.chat_expire", 604800)
	viperV.SetDefault("video.feed_limit", 30)
	viperV.SetDefault("video.popular_limit", 100)
	viperV.SetDefault("expire.refresh_token", 604800)
	viperV.SetDefault("expire.mfa", 604800)
	viperV.SetDefault("expire.chat", 604800)

	if err := viperV.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	Conf = &Config{}
	if err := viperV.Unmarshal(Conf); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 应用环境变量覆盖
	applyEnvOverrides()

	viperV.WatchConfig() // 监听配置文件变更
	viperV.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("配置文件已变更: %s", e.Name)
		reloadConfig() // 重新加载配置
	})

	log.Println("配置初始化完成，支持热更新")
}

func applyEnvOverrides() {
	if dsn := viperV.GetString("MYSQL_DSN"); dsn != "" {
		viperV.Set("mysql.dsn", dsn)
		log.Printf("使用环境变量 MYSQL_DSN: %s", dsn)
	}
	if addr := viperV.GetString("REDIS_ADDR"); addr != "" {
		viperV.Set("redis.addr", addr)
		log.Printf("使用环境变量 REDIS_ADDR: %s", addr)
	}
	if jwtSecret := viperV.GetString("JWT_SECRET"); jwtSecret != "" {
		viperV.Set("jwt.secret", jwtSecret)
		log.Printf("使用环境变量 JWT_SECRET")
	}
}

func reloadConfig() {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	newConf := &Config{}
	if err := viperV.Unmarshal(newConf); err != nil {
		log.Printf("重新加载配置失败: %v", err)
		return
	}

	applyEnvOverrides()

	oldConf := Conf
	Conf = newConf

	log.Printf("配置热更新成功")

	for _, callback := range onChangeCallbacks {
		go callback(newConf)
	}

	_ = oldConf
}

// 注册配置变更回调
func OnConfigChange(callback func(*Config)) {
	onChangeCallbacks = append(onChangeCallbacks, callback)
}

