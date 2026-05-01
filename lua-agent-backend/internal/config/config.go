package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Storage   StorageConfig   `yaml:"storage"`
	LLM       LLMConfig       `yaml:"llm"`
	Embedding EmbeddingConfig `yaml:"embed"`
	Validator ValidatorConfig `yaml:"validator"`
	Logging   LoggingConfig   `yaml:"logging"`
}

type ServerConfig struct {
	Port            string        `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type StorageConfig struct {
	DSN string `yaml:"dsn"`
}

type LLMConfig struct {
	Endpoint   string        `yaml:"endpoint"`
	Model      string        `yaml:"model"`
	Timeout    time.Duration `yaml:"timeout"`
	NumCtx     int           `yaml:"num_ctx"`
	NumPredict int           `yaml:"num_predict"`
	Batch      int           `yaml:"batch"`
	Parallel   int           `yaml:"parallel"`
}

type EmbeddingConfig struct {
	Endpoint string        `yaml:"endpoint"`
	Timeout  time.Duration `yaml:"timeout"`
}

type ValidatorConfig struct {
	Timeout time.Duration `yaml:"timeout"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := defaultConfig()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	return &cfg, nil
}

func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Port:            "8080",
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		LLM: LLMConfig{
			Endpoint:   "http://127.0.0.1:11434",
			Model:      "qwen2.5-coder:7b",
			Timeout:    60 * time.Second,
			NumCtx:     4096,
			NumPredict: 256,
			Batch:      1,
			Parallel:   1,
		},
		Embedding: EmbeddingConfig{
			Endpoint: "http://127.0.0.1:8081",
			Timeout:  30 * time.Second,
		},
		Validator: ValidatorConfig{
			Timeout: 2 * time.Second,
		},
		Logging: LoggingConfig{
			Level: "info",
		},
	}
}
