package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

type conf struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
}

func LoadConfig(path string) (*conf, error) {
	fmt.Println("Tentando carregar configurações do diretório:", path)

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AddConfigPath(filepath.Join(path, ".."))       // Diretório pai
	viper.AddConfigPath(filepath.Join(path, "..", "..")) // Diretório avô
	viper.AddConfigPath("/")                             // Raiz do sistema de arquivos
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
	}

	var config conf
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar configurações: %w", err)
	}

	fmt.Printf("Arquivo de configuração usado: %s\n", viper.ConfigFileUsed())
	fmt.Printf("Configurações carregadas: %+v\n", config)

	return &config, nil
}
