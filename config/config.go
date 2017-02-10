package config

import "github.com/BurntSushi/toml"

type (
	Config struct {
		Title string
		// Server         server
		// Authentication authentication
		MongoDB mongoDB
		NSQ     nsq
		Redis   redis
		JPush   jpush
		Qiniu   qiniu
	}

	server struct {
		Host                  string
		ServerKeyPath         string
		ServerCertificatePath string
	}

	authentication struct {
		PrivateKeyPath string
		PublicKeyPath  string
		TokenDuration  int
		ExpireOffset   int
	}

	mongoDB struct {
		Host string
	}

	nsq struct {
		Host string
	}

	redis struct {
		Host        string
		Maxidle     int
		Maxactive   int
		Idletimeout int
	}

	jpush struct {
		AppKey string
		Secret string
	}

	qiniu struct {
		AccessKey string
		SecretKey string
	}
)

func New() Config {
	var (
		config Config
		err    error
	)
	_, err = toml.DecodeFile("./config/conf.toml", &config)
	if err != nil {
		panic(err)
	}
	return config
}
