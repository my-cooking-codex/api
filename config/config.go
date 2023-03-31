package config

type DBConfig struct {
	URI  string `env:"URI,notEmpty"`
	Type string `env:"TYPE,notEmpty"`
}

type AppConfig struct {
	Host                 string        `env:"HOST" envDefault:"127.0.0.1"`
	Port                 uint          `env:"PORT" envDefault:"8000"`
	DataPath             string        `env:"DATA_PATH,notEmpty"`
	DB                   DBConfig      `envPrefix:"DB__"`
	JWTSecret            Base64Decoded `env:"JWT_SECRET,notEmpty"`
	StaticPath           *string       `env:"STATIC_PATH"`
	CORSOrigins          []string      `env:"CORS_ORIGINS" envSeparator:"," envDefault:"*"`
	OptimizedImageSize   uint          `env:"OPTIMIZED_IMAGE_SIZE" envDefault:"2000"`
	ImageUploadSizeLimit string        `env:"MAX_UPLOAD_SIZE" envDefault:"4M"`
}
