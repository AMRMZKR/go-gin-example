package setting

import (
	"log"
	"time"

	"os"

	"github.com/apsdehal/go-logger"
	"github.com/go-ini/ini"
	"github.com/sakirsensoy/genv"
)

type App struct {
	JwtSecret string
	PageSize  int
	PrefixUrl string

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

var cfg *ini.File

var _log *logger.Logger

// Setup initialize the configuration instance
func Setup() {
	var err error
	_log, err = logger.New("test", 1, os.Stdout)
	if err != nil {
		panic(err) // Check for error
	}
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("redis", RedisSetting)

	initDatabaseConfig(DatabaseSetting)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

func initDatabaseConfig(v *Database) {
	v.Type = genv.Key("DB_TYPE").Default("mysql").String()
	v.Host = genv.Key("DB_HOST").Default("127.0.0.1").String()
	v.User = genv.Key("DB_USER").Default("root").String()
	v.Password = genv.Key("DB_PASSWORD").Default("secret").String()
	v.Name = genv.Key("DB_NAME").Default("test").String()
	v.TablePrefix = genv.Key("DB_TABLE_PREFIX").Default("").String()
	_log.InfoF("%v", *v)
}
