package main

import "fmt"

const (
	RELP_HOSTNAME_OPT = "relpHostname"
	RELP_PORT_OPT     = "relpPort"
	RELP_TLS_OPT      = "relpTls"
	DRIVER_NAME       = "teragrep-relp-docker-plugin"
)

func ValidateLogOpts(cfg map[string]string) error {
	for key := range cfg {
		switch key {
		case RELP_PORT_OPT:
		case RELP_HOSTNAME_OPT:
		case RELP_TLS_OPT:
		default:
			{
				return fmt.Errorf("unknown log opt '%s' provided for %s log driver\n", key, DRIVER_NAME)
			}
		}
	}

	return nil
}
