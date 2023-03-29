package config

type DBConfig struct {
	URI  string `env:"URI,notEmpty"`
	Type string `env:"TYPE,notEmpty"`
}

type AppConfig struct {
	Host       string        `env:"HOST" envDefault:"127.0.0.1"`
	Port       uint          `env:"PORT" envDefault:"8000"`
	DataPath   string        `env:"DATA_PATH,notEmpty"`
	DB         DBConfig      `envPrefix:"DB__"`
	SecretKey  Base64Decoded `env:"SECRET_KEY,notEmpty"`
	StaticPath *string       `env:"STATIC_PATH"`
}
