package config

import "fmt"

type BindConfig struct {
	Host string `env:"HOST" envDefault:"127.0.0.1"`
	Port uint   `env:"PORT" envDefault:"8000"`
}

func (c *BindConfig) AsAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type DBConfig struct {
	URI  string `env:"URI,notEmpty"`
	Type string `env:"TYPE,notEmpty"`
}

type AppConfig struct {
	Bind                 BindConfig    `envPrefix:"BIND__"`
	DB                   DBConfig      `envPrefix:"DB__"`
	DataPath             string        `env:"DATA_PATH,notEmpty"`
	JWTSecret            Base64Decoded `env:"JWT_SECRET,notEmpty"`
	StaticPath           *string       `env:"STATIC_PATH"`
	CORSOrigins          []string      `env:"CORS_ORIGINS" envSeparator:"," envDefault:"*"`
	OptimizedImageSize   uint          `env:"OPTIMIZED_IMAGE_SIZE" envDefault:"2000"`
	ImageUploadSizeLimit string        `env:"MAX_UPLOAD_SIZE" envDefault:"4M"`
}
