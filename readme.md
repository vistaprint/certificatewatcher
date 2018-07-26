
# Certificate Watcher

Certificate Watcher monitors the expiration date of tls certificates stored as secrets in kubernetes and begins alerting when a threshold is crossed

## Running Locally

The following are prerequisites to build and run Certificate Watcher:
* Go 1.7+
* Access to a kubernetes cluster (in-cluster or with a local kubernetes config file)

### Setup Steps

Clone the git repo

Make sure the dependencies are installed
```bash
$ go get github.com/sirupsen/logrus
$ go get gopkg.in/alecthomas/kingpin.v2
$ go get k8s.io/apimachinery/pkg/apis/meta/v1
$ go get k8s.io/client-go/kubernetes
$ go get k8s.io/client-go/tools/clientcmd
$ go install
```

Build the executable

```bash
$ go build
```

Run the application

```bash
$ certificatewatcher --kubeconfig=[PATH TO K8s config] --sender-email=[FROM EMAIL ADDRESS] --notify-email=[TO EMAIL ADDRESS] --smtp-host=[SMTP_HOST] --smtp-port=[SMTP_PORT] --smtp-user=[SMTP_USERNAME] --smtp-password=(SMTP_PASSWORD)
```
### Command Line Arguments

The following are the supported command line arguments for CertificateWatcher

```
--kubeconfig: Path to a kubernetes configuration file (optional)
    * if not specified, default to InCluster config

--interval: interval between two execution cycles (optional) [1m, 1h, etc] 
    * defaults to 1 hour interval

--warning-days: Days from certificate expiration that will trigger a warning (optional)
    * defaults to 30 days

--smtp-host: Host to use when sending warning emails

--smtp-user: Smtp username

--smtp-password: Smtp password

--smtp-port: port to use when sending warning emails
    * defaults to port 25

--notify-email: email to send warning emails (to: address)

--sender-email: email from which to send emails (from: address)

--cluster-name: name of cluster that will appear in subject (optional)

--once: run the loop once (optional)
    * defaults to FALSE
```

### Example Email

```
to: group@example.com
from: alerts@example.com
subject: Warning: Certificates Expiring Soon  - my-cluster

ns=default   cert=example-cert   expires=2020-05-21 19:30:18 +0000 UTC

```