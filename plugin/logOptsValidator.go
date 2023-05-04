package main

import "fmt"

const (
	RELP_HOSTNAME_OPT                      = "relpHostname"
	RELP_PORT_OPT                          = "relpPort"
	RELP_TLS_OPT                           = "relpTls"
	SYSLOG_HOSTNAME_OPT                    = "syslogHostname"
	SYSLOG_APPNAME_OPT                     = "syslogAppName"
	TAG_OPT                                = "tags"
	K8S_METADATA_EXTRACTION_ENABLED_OPT    = "k8sMetadata"
	KUBELET_URL_OPT                        = "kubeletUrl"
	K8S_CLIENT_AUTH_ENABLED_OPT            = "k8sClientAuthEnabled"
	K8S_CLIENT_AUTH_CERT_LOCATION_OPT      = "k8sClientAuthCertPath"
	K8S_CLIENT_AUTH_KEY_LOCATION_OPT       = "k8sClientAuthKeyPath"
	K8S_SERVER_CERT_VALIDATION_ENABLED_OPT = "k8sServerCertValidationEnabled"
	K8S_SERVER_CERT_LOCATION_OPT           = "k8sServerCertPath"
	K8S_METADATA_REFRESH_INTERVAL_OPT      = "k8sMetadataRefreshInterval"
	K8S_METADATA_REFRESH_ENABLED_OPT       = "k8sMetadataRefreshEnabled"
	DRIVER_NAME                            = "teragrep-dkr_01"
)

func ValidateLogOpts(cfg map[string]string) error {
	for key := range cfg {
		switch key {
		case RELP_PORT_OPT:
		case RELP_HOSTNAME_OPT:
		case RELP_TLS_OPT:
		case SYSLOG_HOSTNAME_OPT:
		case SYSLOG_APPNAME_OPT:
		case TAG_OPT:
		case K8S_METADATA_EXTRACTION_ENABLED_OPT:
		case KUBELET_URL_OPT:
		case K8S_CLIENT_AUTH_ENABLED_OPT:
		case K8S_CLIENT_AUTH_CERT_LOCATION_OPT:
		case K8S_CLIENT_AUTH_KEY_LOCATION_OPT:
		case K8S_SERVER_CERT_VALIDATION_ENABLED_OPT:
		case K8S_METADATA_REFRESH_INTERVAL_OPT:
		case K8S_METADATA_REFRESH_ENABLED_OPT:
		default:
			{
				return fmt.Errorf("unknown log opt '%s' provided for %s log driver\n", key, DRIVER_NAME)
			}
		}
	}

	return nil
}
