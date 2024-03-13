package setting

import (
	"gopkg.in/ini.v1"
	"log"
	"time"
)

var cfg *ini.File

// Server settings
type Server struct {
	RunMode         string
	HttpPort        int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	LogSavePath     string
	LogSaveName     string
	LogFileExt      string
	TimeFormat      string
	RuntimeRootPath string
}

var ServerSetting = &Server{}

// JWT settings
type JWT struct {
	AccessExpireTime  time.Duration
	RefreshExpireTime time.Duration
	AccessSecret      string
	RefreshSecret     string
}

var JWTSetting = &JWT{}

// Database settings
type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

var DatabaseSetting = &Database{}

// Redis settings
type Redis struct {
	Host            string
	Password        string
	DB              int
	PoolSize        int
	MinIdleConns    int
	ConnMaxIdleTime time.Duration
}

var RedisSetting = &Redis{}

func Setup() {

	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Printf("not found directory or file at conf/app.ini, so reset directory from test")
		cfg, err = ini.Load("../../conf/app.ini")
		if err != nil {
			log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
		}
	}

	//server settings
	mapTo("server", ServerSetting)
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	//jwt settings
	mapTo("jwt", JWTSetting)
	JWTSetting.AccessExpireTime = JWTSetting.AccessExpireTime * time.Minute
	JWTSetting.RefreshExpireTime = JWTSetting.RefreshExpireTime * time.Hour

	//mariadb settings
	mapTo("database", DatabaseSetting)

	//redis settings
	mapTo("redis", RedisSetting)
	RedisSetting.ConnMaxIdleTime = RedisSetting.ConnMaxIdleTime * time.Second
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
