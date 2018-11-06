package conf

import "github.com/pelletier/go-toml"

type Configuration struct {
	Port          string
	AdminUsername string
	AdminPassword string
	JwtSecret     string
}

func GetConfigurations(path string) (*Configuration, error) {
	conf, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}
	newConf := Configuration{}
	port, ok := conf.Get("port").(string)
	if ok {
		newConf.Port = port
	} else {
		newConf.Port = "8080"
	}

	adminUsername, ok := conf.Get("admin_username").(string)
	if ok {
		newConf.AdminUsername = adminUsername
	} else {
		newConf.AdminUsername = "admin"
	}

	adminPassword, ok := conf.Get("admin_password").(string)
	if ok {
		newConf.AdminPassword = adminPassword
	} else {
		newConf.AdminPassword = "admin"
	}

	jwtSecret, ok := conf.Get("jsw_secret").(string)
	if ok {
		newConf.JwtSecret = jwtSecret
	} else {
		newConf.JwtSecret = "secret"
	}
	return &newConf, nil
}
