package config

import (
	"flag"
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// ErrDatabaseURIRequired возвращается когда не задан обязательный параметр DATABASE_URI.
var ErrDatabaseURIRequired = fmt.Errorf("обязательный параметр DATABASE_URI не задан")

// ErrAccrualSystemAddressRequired возвращается когда не задан обязательный параметр ACCRUAL_SYSTEM_ADDRESS.
var ErrAccrualSystemAddressRequired = fmt.Errorf("обязательный параметр ACCRUAL_SYSTEM_ADDRESS не задан")

// AppConfig представляет конфигурацию приложения.
type AppConfig struct {
	RunAddress           string `envconfig:"RUN_ADDRESS" default:":8080"`
	DatabaseURI          string `envconfig:"DATABASE_URI"`
	AccrualSystemAddress string `envconfig:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecret            string `envconfig:"JWT_SECRET" default:"gophermart-secret-key"`
}

// LoadConfig загружает конфигурацию из переменных окружения и флагов командной строки.
func LoadConfig() (*AppConfig, error) {
	var cfg AppConfig

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения переменных окружения: %w", err)
	}

	var flagRunAddress string
	flag.StringVar(&flagRunAddress, "a", "", "адрес и порт запуска сервиса")

	var flagDatabaseURI string
	flag.StringVar(&flagDatabaseURI, "d", "", "адрес подключения к базе данных")

	var flagAccrualSystemAddress string
	flag.StringVar(&flagAccrualSystemAddress, "r", "", "адрес системы расчёта начислений")

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "a":
			cfg.RunAddress = flagRunAddress
		case "d":
			cfg.DatabaseURI = flagDatabaseURI
		case "r":
			cfg.AccrualSystemAddress = flagAccrualSystemAddress
		}
	})

	if args := flag.Args(); len(args) > 0 {
		return nil, fmt.Errorf("неизвестные аргументы: %v", args)
	}

	if cfg.DatabaseURI == "" {
		return nil, ErrDatabaseURIRequired
	}
	if cfg.AccrualSystemAddress == "" {
		return nil, ErrAccrualSystemAddressRequired
	}

	return &cfg, nil
}
