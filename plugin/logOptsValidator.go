package main

import "fmt"

const (
	RELP_HOSTNAME_OPT = "relpHostname"
	RELP_PORT_OPT     = "relpPort"
	RELP_TLS_OPT      = "relpTls"
	SYSLOG_HOSTNAME   = "syslogHostname"
	SYSLOG_APPNAME    = "syslogAppName"
	DRIVER_NAME       = "teragrep-dkr_01"
)

func ValidateLogOpts(cfg map[string]string) error {
	for key := range cfg {
		switch key {
		case RELP_PORT_OPT:
		case RELP_HOSTNAME_OPT:
		case RELP_TLS_OPT:
		case SYSLOG_HOSTNAME:
		case SYSLOG_APPNAME:
		default:
			{
				return fmt.Errorf("unknown log opt '%s' provided for %s log driver\n", key, DRIVER_NAME)
			}
		}
	}

	return nil
}
