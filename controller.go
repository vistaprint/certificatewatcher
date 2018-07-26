package main

import (
	"crypto/x509"
	"encoding/pem"
	"net/smtp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type certInfo struct {
	Namespace       string
	CertificateName string
	NotAfter        time.Time
}

//Controller type handles the control loop
type Controller struct {
	Config
	KubeClient *kubernetes.Clientset
}

//NewController creates a new controller with the config
func NewController(cfg *Config) (Controller, error) {

	ctrl := Controller{}
	kubeConfig, err := buildConfig(cfg)

	if err != nil {
		return ctrl, err
	}

	ctrl.Config = *cfg
	ctrl.KubeClient = kubeConfig

	return ctrl, nil
}

func buildConfig(cfg *Config) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", cfg.KubeConfig)

	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return client, nil
}

// Run starts the execution loop
func (ctrl Controller) Run(stopChan <-chan struct{}) {

	ticker := time.NewTicker(ctrl.Config.Interval)
	defer ticker.Stop()
	for {
		log.Info("Starting certificate check...")
		err := ctrl.RunOnce()
		if err != nil {
			log.Error(err)
		}
		select {
		case <-ticker.C:
		case <-stopChan:
			log.Error("Terminating main controller loop")
			return
		}
	}
}

// RunOnce executes one iteration of the cert check
func (ctrl Controller) RunOnce() error {
	log.Info("Retrieving tls secrets...")
	nl, err := ctrl.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})

	if err != nil {
		return err
	}

	expired := []certInfo{}
	threshold := time.Now().AddDate(0, 0, ctrl.Config.WarningDays)
	for _, n := range nl.Items {
		namespace := n.ObjectMeta.Name

		sl, err := ctrl.KubeClient.CoreV1().Secrets(namespace).List(metav1.ListOptions{
			FieldSelector: "type=kubernetes.io/tls",
		})

		if err != nil {
			return err
		}

		for _, sec := range sl.Items {
			block, _ := pem.Decode(sec.Data["tls.crt"])
			cert, err := x509.ParseCertificate(block.Bytes)

			if err != nil {
				return err
			}

			if cert.NotAfter.Before(threshold) {

				ec := certInfo{
					Namespace:       namespace,
					CertificateName: sec.ObjectMeta.Name,
					NotAfter:        cert.NotAfter}

				expired = append(expired, ec)
			}

		}
	}

	//If there's no expiring certs, return
	if len(expired) == 0 {
		log.Info("No cert alerts detected")
		return nil
	}

	log.Info("Sending notifications")
	err = ctrl.notify(expired)

	if err != nil {
		return err
	}

	return nil
}

func (ctrl Controller) notify(expiredInfo []certInfo) error {
	auth := smtp.PlainAuth("", ctrl.Config.SMTPUsername, ctrl.Config.SMTPPassword, ctrl.Config.SMTPHost)

	to := []string{ctrl.Config.NotifyEmailAddr}

	// Build the message
	msgStr := "To: " + ctrl.Config.NotifyEmailAddr + "\r\n" +
		"Subject: Warning: Certificates Expiring Soon " + " - " + ctrl.Config.ClusterName + "\r\n" +
		"\r\n"

	for _, ci := range expiredInfo {
		log.Infof("%v", ci)
		msgStr += "ns=" + ci.Namespace + "   cert=" + ci.CertificateName + "   expires=" + ci.NotAfter.String() + "\r\n"
	}

	log.Info("sending alert")
	err := smtp.SendMail(ctrl.SMTPHost+":"+strconv.Itoa(ctrl.SMTPPort), auth, ctrl.Config.SenderEmailAddr, to, []byte(msgStr))

	if err != nil {
		return err
	}

	return nil
}
