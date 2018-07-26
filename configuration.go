package main

import (
	"strconv"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

type Config struct {
	Interval        time.Duration
	WarningDays     int
	KubeConfig      string
	SMTPHost        string
	SMTPUsername    string
	SMTPPassword    string
	SMTPPort        int
	NotifyEmailAddr string
	ClusterName     string
	Once            bool
	SenderEmailAddr string
}

var defConfig = &Config{
	Interval:        time.Hour,
	KubeConfig:      "",
	WarningDays:     30,
	SMTPHost:        "",
	SMTPUsername:    "",
	SMTPPassword:    "",
	SMTPPort:        25,
	NotifyEmailAddr: "",
	ClusterName:     "",
	Once:            false,
	SenderEmailAddr: "",
}

// NewConfig creates a new default config
func NewConfig() *Config {
	return &Config{}
}

// ParseFlags parses arguments into a config
func (cfg *Config) ParseFlags(args []string) error {
	app := kingpin.New("certificate-watcher", "watch tls secrets in a k8s cluster")

	// this should come from the build number
	app.Version("1.0.0")

	app.DefaultEnvars()
	app.Flag("kubeconfig", "K8s config file to use").Default(defConfig.KubeConfig).StringVar(&cfg.KubeConfig)
	app.Flag("interval", "Interval between two execution cycles (default 1 month)").Default(defConfig.Interval.String()).DurationVar(&cfg.Interval)
	app.Flag("warning-days", "Days from certificate expiration date that will trigger a warning notification (default 30)").Default(strconv.Itoa(defConfig.WarningDays)).IntVar(&cfg.WarningDays)
	app.Flag("smtp-host", "smtp host for email notifications (optional)").Default(defConfig.SMTPHost).StringVar(&cfg.SMTPHost)
	app.Flag("smtp-user", "smtp username for email notifications (optional)").Default(defConfig.SMTPUsername).StringVar(&cfg.SMTPUsername)
	app.Flag("smtp-password", "smtp password for email notifications (optional)").Default(defConfig.SMTPPassword).StringVar(&cfg.SMTPPassword)
	app.Flag("smtp-port", "smtp port to use for notifications (default 25)").Default(strconv.Itoa(defConfig.SMTPPort)).IntVar(&cfg.SMTPPort)
	app.Flag("notify-email", "email address to send cert expiration warnings").Default(defConfig.NotifyEmailAddr).StringVar(&cfg.NotifyEmailAddr)
	app.Flag("cluster-name", "name of cluster to use in notifications (optional)").Default(defConfig.ClusterName).StringVar(&cfg.ClusterName)
	app.Flag("once", "When enabled, the control loop is broken after one run (default: disabled)").BoolVar(&cfg.Once)
	app.Flag("sender-email", "email address from which to send cert expiration warnings").Default(defConfig.SenderEmailAddr).StringVar(&cfg.SenderEmailAddr)
	_, err := app.Parse(args)

	if err != nil {
		return err
	}

	return nil
}
