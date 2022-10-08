package Config

type MysqlConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	Name     string `json:"db" mapstructure:"db"`
	User     string `json:"user" mapstructure:"user"`
	PassWord string `json:"password" mapstructure:"password"`
}

type CosulConfig struct {
	Host string `json:"host" mapstructure:"host"`
	Port int    `json:"port" mapstructure:"port"`
}
type ServerConfig struct {
	Host            string      `json:"host" mapstructure:"host"`
	Port            int         `json:"port" mapstructure:"port"`
	Name            string      `json:"name" mapstructure:"name"`
	MysqlInfo       MysqlConfig `mapstructure:"mysql" json:"mysql"`
	CosulConfigInfo CosulConfig `mapstructure:"consul" json:"consul"`
}
