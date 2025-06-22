package util

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

type SecretManager struct {
	cache     *secretcache.Cache
	secretMap map[string]string
	mu        sync.Once
}

var (
	secretManager     *SecretManager
	secretManagerOnce sync.Once
)

func NewSecretManager() (*SecretManager, error) {
	cache, err := secretcache.New()

	if err != nil {
		return nil, err
	}

	return &SecretManager{cache: cache}, nil
}

func (sm *SecretManager) Load() error {
	var err error

	sm.mu.Do(func() {
		secretName := os.Getenv("SECRET_MANAGER_NAME")
		if secretName == "" {
			err = fmt.Errorf("SECRET_MANAGER_NAME is not set")
			return
		}

		secretString, cacheErr := sm.cache.GetSecretStringWithContext(context.Background(), secretName)
		if cacheErr != nil {
			err = fmt.Errorf("failed to retrieve secret -> %v", cacheErr)
			return
		}

		secretData := make(map[string]string)
		if unmarshalErr := json.Unmarshal([]byte(secretString), &secretData); unmarshalErr != nil {
			err = fmt.Errorf("failed to parse secret JSON \n JSON: %v", unmarshalErr)
			return
		}

		sm.secretMap = secretData
	})

	return err
}

func (sm *SecretManager) Get(secretKey string) (*string, error) {
	if sm.secretMap == nil {
		if err := sm.Load(); err != nil {
			return nil, err
		}
	}

	value, exists := sm.secretMap[secretKey]
	if !exists {
		return nil, fmt.Errorf("secret key not found %s", secretKey)
	}

	return &value, nil
}

func GetSecretValue(secretName string) string {
	env := os.Getenv("ENV")

	if env != "LOCAL" {
		secretManagerOnce.Do(func() {
			var err error

			secretManager, err = NewSecretManager()

			if err != nil {
				Infof("⚠️  SecretManager initialization failed: %v. Falling back to .env file...", err)
			}
		})

		if secretManager != nil {
			if value, err := secretManager.Get(secretName); err == nil {
				return *value
			} else {
				Infof("⚠️  Secret '%s' not found in SecretManager. Falling back to environment variable: %v", secretName, err)
			}
		}
	} else {
		GetLogger().Info("Using local environment, use environment variables instead of secret manager.")
	}

	if err := godotenv.Load(".env"); err != nil {
		Errorf("❌ Unable to load environment variables from .env file: %v", err)
	}

	envValue := os.Getenv(secretName)

	if envValue == "" {
		Errorf("⚠️  Secret '%s' is missing from both SecretManager and environment variables.", secretName)
	}

	return envValue
}
