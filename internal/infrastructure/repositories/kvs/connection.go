package kvs

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mercadolibre/fury_go-toolkit-config/pkg/config"
	"github.com/mercadolibre/go-meli-toolkit/gokvsclient"
)

const productionEnv string = "production"

var (
	environment   = map[string]string{}
	maxIdleConn   = "KVS_MAX_IDLE_CONN"
	timeout       = "KVS_TIMEOUT"
	containerName = "KVS_CONTAINER_NAME"
	defaultString = ""
)

func GetKVSConnection() (gokvsclient.Client, error) {
	env := os.Getenv("GO_ENVIRONMENT")
	if env == productionEnv {
		if cfg, err := config.Load(); err != nil {
			return nil, err
		} else {
			environment = map[string]string{
				containerName: cfg.GetString(containerName, defaultString),
				maxIdleConn:   cfg.GetString(maxIdleConn, defaultString),
				timeout:       cfg.GetString(timeout, defaultString),
			}
		}
	} else {
		loadENVConfigs()
	}

	kvsConfig := gokvsclient.MakeKvsConfig()
	kvsConfig.SetReadMaxIdleConnections(50)
	kvsConfig.SetWriteMaxIdleConnections(50)
	kvsConfig.SetReadTimeout(150 * time.Millisecond)
	kvsConfig.SetWriteTimeout(150 * time.Millisecond)

	return gokvsclient.MakeKvsClient(environment[containerName], kvsConfig), nil
}

func loadENVConfigs() {
	if err := godotenv.Load("../../.env"); err != nil {
		panic(err.Error())
	}

	environment = map[string]string{
		containerName: os.Getenv(containerName),
		maxIdleConn:   os.Getenv(maxIdleConn),
		timeout:       os.Getenv(timeout),
	}
}
