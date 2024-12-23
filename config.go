package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const defaultAWSRegion = "eu-central-1"

type Config struct {
	IsDev          bool
	TelegramToken  string
	TelegramChatID string
}

func GetConfig() (*Config, error) {
	if os.Getenv("ENV") == "dev" {
		return &Config{
			IsDev:          true,
			TelegramToken:  os.Getenv("TELEGRAM_TOKEN"),
			TelegramChatID: os.Getenv("TELEGRAM_CHAT_ID"),
		}, nil
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultAWSRegion
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("create aws session: %w", err)
	}

	awsConfig := aws.NewConfig().WithRegion(region)
	if os.Getenv("AWS_ENV_CREDS") == "true" {
		awsConfig = awsConfig.WithCredentials(credentials.NewEnvCredentials())
	}
	ssmClient := ssm.New(sess, awsConfig)
	parameters, err := ssmClient.GetParameters(&ssm.GetParametersInput{
		Names: []*string{
			aws.String("/certification-preparation-bot/prod/openai-api-token"),
			aws.String("/certification-preparation-bot/prod/telegram-token"),
			aws.String("/certification-preparation-bot/prod/telegram-chat-id"),
		},
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("get parameters: %w", err)
	}

	openaiAPIToken := ""
	telegramToken := ""
	telegramChatID := ""
	for _, param := range parameters.Parameters {
		switch *param.Name {
		case "/certification-preparation-bot/prod/openai-api-token":
			openaiAPIToken = *param.Value
		case "/certification-preparation-bot/prod/telegram-token":
			telegramToken = *param.Value
		case "/certification-preparation-bot/prod/telegram-chat-id":
			telegramChatID = *param.Value
		}
	}

	errs := make([]string, 0, 3)
	if openaiAPIToken == "" {
		errs = append(errs, "missing openai token")
	}
	if telegramToken == "" {
		errs = append(errs, "missing telegram token")
	}
	if telegramChatID == "" {
		errs = append(errs, "missing telegram chat id")
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("missing required parameters: %s", errs)
	}

	return &Config{
		TelegramToken:  telegramToken,
		TelegramChatID: telegramChatID,
	}, nil
}
