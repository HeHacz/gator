package config

func (conf *Config) SetUser(username string) error {
	conf.Current_user_name = username
	return write(*conf)
}
