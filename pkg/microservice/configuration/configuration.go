package configuration

import (
	"os"
	"path/filepath"
	"time"

	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/utils"
	"github.com/spf13/viper"
)

type Configurationobject struct {
	configurationutl *viper.Viper
}

func (c *Configurationobject) InitGlobalConfiguration(servicename string) (err error) {
	c.configurationutl = viper.New()

	err = utils.EnsurePathExist(c.getRootConfigPath(servicename))
	if err != nil {
		logger.WarnErrLog(err)
	}

	err = c.registerConfigFile(servicename)
	if err != nil {
		return err
	}

	err = c.readConfigFile()
	if err != nil {
		return err
	}

	return nil
}

func (c *Configurationobject) registerConfigFile(servicename string) error {
	var exist error

	const filename = "config.toml"

	rootpath := c.getRootConfigPath(servicename)
	userpath, err := c.getUserConfigPath(servicename)
	if err != nil {
		return err
	}

	c.configurationutl.SetConfigName("config")
	c.configurationutl.SetConfigType("toml")

	// ignore permission errors
	wdpath, getabserr := filepath.Abs(".")
	if getabserr != nil {
		logger.WarnErrLog(getabserr)
	}

	insidecmdpath := filepath.Join("../../configs", servicename)
	_, exist = os.Stat(filepath.Join(insidecmdpath, filename))

	if exist == nil {
		c.configurationutl.AddConfigPath(insidecmdpath)
	}

	insiderootpath := filepath.Join("./configs", servicename)
	_, exist = os.Stat(filepath.Join(insiderootpath, filename))

	if exist == nil {
		c.configurationutl.AddConfigPath(insiderootpath)
	}

	_, exist = os.Stat(filepath.Join(wdpath, filename))
	if exist == nil {
		c.configurationutl.AddConfigPath(wdpath)
	}

	_, exist = os.Stat(filepath.Join(rootpath, filename))
	if exist == nil {
		c.configurationutl.AddConfigPath(rootpath)
	}

	_, exist = os.Stat(filepath.Join(userpath, filename))
	if exist == nil {
		c.configurationutl.AddConfigPath(userpath)
	}

	return nil
}

func (c *Configurationobject) getUserConfigPath(servicename string) (string, error) {
	userConfigDir, err := c.getUserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userConfigDir, servicename), nil
}

func (c *Configurationobject) getRootConfigPath(servicename string) string {
	return filepath.Join("/etc/", servicename)
}

func (c *Configurationobject) getUserConfigDir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return userConfigDir, nil
}

func (c *Configurationobject) readConfigFile() error {
	err := c.configurationutl.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func (c *Configurationobject) GetAllSettings() map[string]interface{} {
	return c.configurationutl.AllSettings()
}

func (c *Configurationobject) GetStringSetting(key string) string {
	return c.configurationutl.GetString(key)
}

func (c *Configurationobject) GetIntSetting(key string) int {
	return c.configurationutl.GetInt(key)
}

func (c *Configurationobject) GetStringSliceSetting(key string) []string {
	return c.configurationutl.GetStringSlice(key)
}

func (c *Configurationobject) GetStringMapSettings(key string) map[string]string {
	return c.configurationutl.GetStringMapString(key)
}

func (c *Configurationobject) GetBoolSetting(key string) bool {
	return c.configurationutl.GetBool(key)
}

func (c *Configurationobject) GetTimeDuration(key string) time.Duration {
	return c.configurationutl.GetDuration(key)
}
