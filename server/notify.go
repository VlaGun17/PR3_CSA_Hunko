package server

import (
	"fmt"
	"os"
	"time"

	"PR3_CSA_Hunko/smtp"
)

const DutyEngineerEmail = "duty-engineer@internal.net"

func smtpConfigFromEnv() smtp.Config {
	return smtp.Config{
		Host:     getenv("SMTP_HOST", "smtp.gmail.com"),
		Port:     getenv("SMTP_PORT", "587"),
		From:     os.Getenv("SMTP_FROM"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

func SendIncidentNotification(incident, clientID string, registeredAt time.Time) error {
	cfg := smtpConfigFromEnv()
	if cfg.From == "" || cfg.Password == "" {
		return fmt.Errorf("SMTP_FROM / SMTP_PASSWORD are not set")
	}

	subject := "[INCIDENT] Нове повідомлення про інцидент"
	body := fmt.Sprintf(
		"Шановний черговий інженере!\r\n\r\n"+
			"Сервер моніторингу отримав нове повідомлення про інцидент через захищене mTLS-з'єднання.\r\n\r\n"+
			"Дата/час реєстрації: %s\r\n"+
			"Джерело (сертифікат клієнта): %s\r\n"+
			"Повідомлення: %s\r\n\r\n"+
			"Це повідомлення згенеровано автоматично.\r\n",
		registeredAt.Format("02.01.2006 15:04:05 MST"),
		clientID,
		incident,
	)

	return smtp.NewSender(cfg).SendEmail([]string{DutyEngineerEmail}, subject, body)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
