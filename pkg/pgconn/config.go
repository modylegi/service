package pgconn

type Config struct {
	host     string
	port     int
	database string
	username string
	password string
}

func NewConfig() *Config {
	return &Config{
		host:     "localhost",
		port:     5432,
		database: "database",
		username: "postgres",
		password: "password",
	}
}

func (c *Config) WithHost(host string) *Config {
	c.host = host
	return c
}

func (c *Config) WithPort(port int) *Config {
	c.port = port
	return c
}

func (c *Config) WithDatabase(database string) *Config {
	c.database = database
	return c
}

func (c *Config) WithUsername(username string) *Config {
	c.username = username
	return c
}

func (c *Config) WithPassword(password string) *Config {
	c.password = password
	return c
}
