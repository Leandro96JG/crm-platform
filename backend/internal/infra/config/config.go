package config

import (
	"os"
	"time"

	"github.com/magiconair/properties"
)

const (
	defaultAppPropertyFilename = "resources/application.properties"
	appPropertyFilenameEnv     = "APP_PROPERTY_FILE"
)

type AppConfig struct {
	Database          Database   `properties:"database"`
	SecretJWTKey      string     `properties:"jwtKeyEnv"`
	ReportFolder      string     `properties:"reportFolder,default=resources/reports"`
	AttachmentsBucket Bucket     `properties:"attachmentBucket"`
	WhatsApp          WhatsApp   `properties:"whatsapp"`
	OpenAI            OpenAIConf `properties:"openai"`
}

type Database struct {
	ConnStr         string        `properties:"connStr,default="`
	Driver          string        `properties:"driver,default=postgres"`
	Port            int           `properties:"port,default=5432"`
	Schema          string        `properties:"schema"`
	Username        string        `properties:"username"`
	PasswordEnv     string        `properties:"passwordEnv"`
	HostEnv         string        `properties:"hostEnv"`
	ServerTimezone  string        `properties:"serverTimezone,default=America/Sao_Paulo"`
	MaxOpenConns    int           `properties:"maxOpenConns,default=50"`
	MaxIdleConns    int           `properties:"maxIdleConns,default=50"`
	ConnMaxLifetime time.Duration `properties:"connMaxLifetime,default=10m"`
}

type Bucket struct {
	Name            string        `properties:"name"`
	Region          string        `properties:"region"`
	Timeout         time.Duration `properties:"timeout"`
	AWSKeyIDEnv     string        `properties:"awsKeyIdEnv"`
	AWSSecretKeyEnv string        `properties:"awsSecretKeyEnv"`
}

func (db Database) Host() string {
	return os.Getenv(db.HostEnv)
}

func (db Database) Password() string {
	return os.Getenv(db.PasswordEnv)
}

func (db AppConfig) SecretKey() string {
	return os.Getenv(db.SecretJWTKey)
}

func (b Bucket) AWSKeyID() string {
	return os.Getenv(b.AWSKeyIDEnv)
}

func (b Bucket) AWSSecretKey() string {
	return os.Getenv(b.AWSSecretKeyEnv)
}

type WhatsApp struct {
	PhoneNumberID string `properties:"phoneNumberId,default="`
	AccessToken   string `properties:"accessTokenEnv,default="`
	APIVersion    string `properties:"apiVersion,default=v21.0"`
	WebhookSecret string `properties:"webhookSecretEnv,default="`
}

func (w WhatsApp) AccessTokenValue() string {
	return os.Getenv(w.AccessToken)
}

func (w WhatsApp) WebhookSecretValue() string {
	return os.Getenv(w.WebhookSecret)
}

type OpenAIConf struct {
	APIKey string `properties:"apiKeyEnv,default="`
	Model  string `properties:"model,default=gpt-4o-mini"`
}

func (o OpenAIConf) APIKeyValue() string {
	return os.Getenv(o.APIKey)
}

func AppPropertyFilename() string {
	propFile := defaultAppPropertyFilename
	if f := os.Getenv(appPropertyFilenameEnv); f != "" {
		propFile = f
	}
	return propFile
}

func Load() (*AppConfig, error) {
	p := properties.MustLoadFile(AppPropertyFilename(), properties.UTF8)

	config := &AppConfig{}
	err := p.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
