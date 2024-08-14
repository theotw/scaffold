package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"

	"github.com/theotw/scaffold/errorutils"
	"github.com/theotw/scaffold/logger"
)

// wellKnownKeys is a map of environment variables that are well known and should included in the settings
// even though they do not have the prefix
var wellKnownKeys = map[string]bool{
	"AWS_ACCESS_KEY_ID":     true,
	"AWS_SECRET_ACCESS_KEY": true,
	"AWS_PROFILE":           true,
	"AWS_DEFAULT_REGION":    true,
}

type Source[t any] interface {
	InitSettings() error
	GetSettings() *t
	GetItem(key string) string
}

type EnvSettingsSource[t any] struct {
	prefix        string
	fileName      string
	mapOfSettings map[string]string
	settings      *t
}

func NewEnvSettingsSource[t any](prefix string, fileName string) *EnvSettingsSource[t] {
	return &EnvSettingsSource[t]{prefix: prefix, fileName: fileName}
}
func (f *EnvSettingsSource[t]) InitSettings() error {
	log := logger.NewLogger(context.Background(), "settings")

	if err := godotenv.Load(f.fileName); err != nil {
		return errorutils.WrapIfNotNil("error loading settings file "+f.fileName, err)
	} else {
		log.Debugf("Loaded settings from %s", f.fileName)
	}

	mapFromEnv := f.readSettingsFromEnv()
	awsSecretData, err := f.loadAwsSecret()
	if err != nil {
		return errorutils.WrapIfNotNil("error loading aws secret", err)
	}
	if awsSecretData != nil {
		f.mapOfSettings = awsSecretData
		// let the env override aws secrets
		for k, v := range mapFromEnv {
			f.mapOfSettings[k] = v
		}
	} else {
		log.Debugf("No AWS secret loaded")
		f.mapOfSettings = mapFromEnv
	}

	bits, err := json.Marshal(f.mapOfSettings)
	if err != nil {
		return errorutils.WrapIfNotNil("error marshalling settings", err)
	}
	err = json.Unmarshal(bits, &f.settings)
	if err != nil {
		return errorutils.WrapIfNotNil("error unmarshalling settings", err)
	}

	return nil
}
func (f *EnvSettingsSource[t]) GetSettings() *t {
	return f.settings
}
func (f *EnvSettingsSource[t]) GetItem(key string) string {
	return f.mapOfSettings[key]
}

func (f *EnvSettingsSource[t]) readSettingsFromEnv() map[string]string {
	envVars := os.Environ()
	mapOfSettings := make(map[string]string)
	for _, envVar := range envVars {
		parts := strings.Split(envVar, "=")
		if len(parts) != 2 {
			continue
		}
		envKey := parts[0]
		value := parts[1]
		if strings.HasPrefix(envKey, f.prefix) {
			key := envKey[len(f.prefix):]
			key = strings.ToLower(key)
			mapOfSettings[key] = value
		} else {
			if _, ok := wellKnownKeys[envKey]; ok {
				key := strings.ToLower(envKey)
				mapOfSettings[key] = value
			}
		}
	}
	return mapOfSettings
}

func (f *EnvSettingsSource[t]) loadAwsSecret() (map[string]string, error) {
	secretName := os.Getenv(f.prefix + "AWS_SECRET_NAME")
	if secretName == "" {
		return nil, nil
	}
	awsSecretData := make(map[string]string)
	//pull these out first, because we need them to load settings
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsProfile := os.Getenv("AWS_PROFILE")

	cfg, err := loadConfigWithKeysOrProfile(awsAccessKey, awsSecretKey, awsProfile)
	if err != nil {
		return nil, err
	}

	secretsManagerClient := secretsmanager.NewFromConfig(*cfg)
	secretValue, err := secretsManagerClient.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		return nil, errorutils.WrapIfNotNil(fmt.Sprintf("error loading aws secret %s", secretName), err)
	}
	if err := json.Unmarshal([]byte(*secretValue.SecretString), &awsSecretData); err != nil {
		return nil, errorutils.WrapIfNotNil("error unmarshalling secret data", err)
	}

	return awsSecretData, nil
}

// loadConfigWithKeysOrProfile loads the AWS config with the provided access keys or profile.
func loadConfigWithKeysOrProfile(accessKeyID, secretAccessKey, profile string) (*aws.Config, error) {
	ctx := context.Background()
	log := logger.NewLogger(context.Background(), "settings")
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		log.Infof("No region set, defaulting to us-east-1")
		region = "us-east-1"
	}
	var cfg aws.Config
	var err error
	if accessKeyID != "" && secretAccessKey != "" {
		log.Debugf("Using access keys from Environment")
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")))
	} else {
		log.Debugf("No access key set, Using profile %s\n", profile)
		if profile == "" {
			log.Debugf("No profile set, using default")
			cfg, err = config.LoadDefaultConfig(ctx)
		} else {
			cfg, err = config.LoadDefaultConfig(ctx,
				config.WithSharedConfigProfile(profile))
		}
	}
	returnErr := errorutils.WrapIfNotNil("error loading AWS config", err)
	return &cfg, returnErr

}
