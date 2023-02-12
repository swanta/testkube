package config

import (
	"encoding/json"
	"os"
	"path"

	"github.com/kubeshop/testkube/pkg/oauth"
)

var DefaultConfig = Data{
	TelemetryEnabled: true,
	Namespace:        "testkube",
	APIURI:           "http://localhost:8088",
	OAuth2Data: OAuth2Data{
		Provider: oauth.GithubProviderType,
	},
}

const configDirName = ".testkube"
const configFile = "config.json"

func GetStorage(dir string) (Storage, error) {
	storage := Storage{Dir: dir}
	err := storage.Init()
	return storage, err
}

type Storage struct {
	Dir string
}

func (c *Storage) Load() (data Data, err error) {
	path, err := c.getPath()
	if err != nil {
		return data, err
	}
	d, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(d, &data)
	return data, err
}

func (c *Storage) Save(data Data) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	path, err := c.getPath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, d, 0700)
}

func (c *Storage) Init() error {
	var defaultConfig = DefaultConfig
	// create ConfigWriter dir if not exists
	dir, err := c.getDir()
	if err != nil {
		return err
	}
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.Mkdir(dir, 0700)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// create empty JSON file if not exists
	path, err := c.getPath()
	if err != nil {
		return err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		f.Close()

		return c.Save(defaultConfig)
	} else if err != nil {
		return err
	}

	return nil
}

func (c *Storage) getDir() (string, error) {
	var err error
	var dir string
	if c.Dir != "" {
		dir = c.Dir
	} else {
		dir, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}

	return path.Join(dir, configDirName), nil
}

func (c *Storage) getPath() (string, error) {
	dir, err := c.getDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, configFile), nil
}
