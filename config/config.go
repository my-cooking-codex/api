package config

import (
	"fmt"
	"path"
)

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

type DataConfig struct {
	RecipeImagesBase string `env:"RECIPE_IMAGES_BASE,notEmpty"`
}

func (c *DataConfig) RecipeOriginalsPath() string {
	return path.Join(c.RecipeImagesBase, "original")
}

type AppConfig struct {
	Bind                 BindConfig    `envPrefix:"BIND__"`
	DB                   DBConfig      `envPrefix:"DB__"`
	Data                 DataConfig    `envPrefix:"DATA__"`
	JWTSecret            Base64Decoded `env:"JWT_SECRET,notEmpty"`
	StaticPath           *string       `env:"STATIC_PATH"`
	CORSOrigins          []string      `env:"CORS_ORIGINS" envSeparator:"," envDefault:"*"`
	OptimizedImageSize   uint          `env:"OPTIMIZED_IMAGE_SIZE" envDefault:"2000"`
	ImageUploadSizeLimit string        `env:"MAX_UPLOAD_SIZE" envDefault:"4M"`
}
