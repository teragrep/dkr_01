package main

import "fmt"

const (
	RELP_HOSTNAME_OPT   = "relpHostname"
	RELP_PORT_OPT       = "relpPort"
	RELP_TLS_OPT        = "relpTls"
	SYSLOG_HOSTNAME_OPT = "syslogHostname"
	SYSLOG_APPNAME_OPT  = "syslogAppName"
	DRIVER_NAME         = "teragrep-dkr_01"

	TAG_OPT = "tag"
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
		default:
			{
				return fmt.Errorf("unknown log opt '%s' provided for %s log driver\n", key, DRIVER_NAME)
			}
		}
	}

	return nil
}
