package config

import "gopkg.in/ini.v1"

// InitSetting
//
//	@Description: 初始化配置
//	@param file 配置文件路径
//	@return error
func InitSetting(file string) (*AppConfig, error) {
	cfg, err := ini.Load(file)
	if err != nil {
		return nil, err
	}
	config := new(AppConfig)
	_ = cfg.MapTo(config)

	config.MySqlConfig = NewMySqlConfig(cfg)
	config.RedisConfig = NewRedisConfig(cfg)
	config.EmailConfig = NewEmailConfig(cfg)
	config.COSConfig = NewCOSConfig(cfg)
	config.FilePathConfig = NewFilePathConfig(cfg)
	return config, nil
}

// AppConfig
// @Description:应用配置
type AppConfig struct {
	Release         bool   `ini:"release"` //是否是上线模式
	Port            string `ini:"port"`    //端口
	ProUrl          string `ini:"proUrl"`
	DefaultPassword string `ini:"defaultPassword"`
	*MySqlConfig
	*RedisConfig
	*EmailConfig
	*ReleasePathConfig
	*COSConfig
	*FilePathConfig
}

type ReleasePathConfig struct {
	StartWith []string
}

// MySqlConfig
// @Description: mysql相关配置
type MySqlConfig struct {
	User     string `ini:"user"`     //用户名
	Password string `ini:"password"` //密码
	DB       string `ini:"db"`       //要操作的数据库
	Host     string `ini:"host"`     //host
	Port     string `ini:"port"`     //端口
}

func NewMySqlConfig(cfg *ini.File) *MySqlConfig {
	mysqlConfig := &MySqlConfig{}
	cfg.Section("mysql").MapTo(mysqlConfig)
	return mysqlConfig
}

// RedisConfig
// @Description: redis相关配置
type RedisConfig struct {
	Host     string `ini:"host"`
	Port     string `ini:"port"`
	Password string `ini:"password"`
}

func NewRedisConfig(cfg *ini.File) *RedisConfig {
	redisConfig := &RedisConfig{}
	cfg.Section("redis").MapTo(redisConfig)
	return redisConfig
}

// EmailConfig
// @Description: email相关配置
type EmailConfig struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	User     string `ini:"user"`
	Password string `ini:"pass"`
}

func NewEmailConfig(cfg *ini.File) *EmailConfig {
	emailConfig := &EmailConfig{}
	cfg.Section("email").MapTo(emailConfig)
	return emailConfig
}

// FilePathConfig
// @Description: 文件路径相关配置
type FilePathConfig struct {
	ProblemFileDir             string `ini:"problemFileDir"`             //题目文件目录
	ProblemDescriptionTemplate string `ini:"problemDescriptionTemplate"` //题目描述模板文位置
	ProblemFileTemplate        string `ini:"problemFileTemplate"`        //题目编程文件的模板文件
	TempDir                    string `ini:"tmpDir"`                     //临时目录

}

func NewFilePathConfig(cfg *ini.File) *FilePathConfig {
	filePathConfig := &FilePathConfig{}
	cfg.Section("filePath").MapTo(filePathConfig)
	return filePathConfig
}

// COSConfig
// @Description:oss相关配置
type COSConfig struct {
	AppID             string `ini:"appID"`
	Region            string `ini:"region"`
	SecretID          string `ini:"secretID"`
	SecretKey         string `ini:"secretKey"`
	ProblemBucketName string `ini:"problemBucketName"`
	ImageBucketName   string `ini:"imageBucketName"`
}

func NewCOSConfig(cfg *ini.File) *COSConfig {
	cosConfig := &COSConfig{}
	cfg.Section("cos").MapTo(cosConfig)
	return cosConfig
}
