package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/y0ssar1an/q"
	"log"
	"ad_exporter/vars"
)

type adMetrics struct {
	account_lockout_status_gauge *prometheus.Desc
	account_disabled_status_gauge *prometheus.Desc
	password_expired_status_gauge *prometheus.Desc
	credential_validation_status_gauge *prometheus.Desc
	unknown_authentication_error_status_gauge *prometheus.Desc
	ldap_connection_status_gauge *prometheus.Desc


}


func ADmetricsCollector() *adMetrics{

return &adMetrics{

	account_lockout_status_gauge: prometheus.NewDesc("account_lockout_status_gauge","This will return 1 if an account is locked in AD or 0",nil,nil,
		),

	account_disabled_status_gauge: prometheus.NewDesc("account_disabled_status_gauge","This will return 1 if an account is disabled in AD or 0",nil,nil,
	),

	password_expired_status_gauge: prometheus.NewDesc("password_expired_status_gauge","This will return 1 if an account password is expired in AD or 0",nil,nil,
	),

	credential_validation_status_gauge: prometheus.NewDesc("credential_validation_status_gauge","This will return 1 if an account credential you provided are different in AD or 0",nil,nil,
	),

	unknown_authentication_error_status_gauge: prometheus.NewDesc("unknown_authentication_error_status_gauge","This will return 1 if authentication failed due to unknown reason",nil,nil,
	),

	ldap_connection_status_gauge: prometheus.NewDesc("ldap_connection_status_gauge","This will return 1 if LDAP connection failed 0 otherwise",nil,nil,
	),



}


}

func (collector *adMetrics) Describe(ch chan<- *prometheus.Desc){


	ch <- collector.account_lockout_status_gauge
	ch <- collector.account_disabled_status_gauge
	ch <- collector.ldap_connection_status_gauge
	ch <- collector.password_expired_status_gauge
	ch <- collector.unknown_authentication_error_status_gauge
	ch <- collector.credential_validation_status_gauge



}



func (collector *adMetrics) Collect(ch chan<- prometheus.Metric) {

	connection := connectToLdap(vars.Inputs.LdapServer, vars.Inputs.LdapServerPort)
	q.Q(connection)
	if connection.error != nil {
		ldapConnectionError = 1

	} else {
		ldapConnectionError = 0
		log.Printf("LDAP connection successfull , going ghead with binding")
		BindError := bindToLdap(connection.connection, vars.Inputs.UserName, vars.Inputs.LdapBindPassword)

		userAccountUnknownError, userAccountPasswordExpired, userAccountDisabled, userAccountLocked, userAccountWrongPassword := analyzeLdapError(BindError)
		log.Printf("Values are ", userAccountUnknownError, userAccountPasswordExpired, userAccountDisabled, userAccountLocked, userAccountWrongPassword)

	}
ch <- prometheus.MustNewConstMetric(collector.account_lockout_status_gauge,prometheus.GaugeValue,vars.Inputs.UserAccountLocked)
ch <- prometheus.MustNewConstMetric(collector.account_disabled_status_gauge,prometheus.GaugeValue,vars.Inputs.UserAccountDisabled)
ch <- prometheus.MustNewConstMetric(collector.ldap_connection_status_gauge,prometheus.GaugeValue,vars.Inputs.LdapConnectionStatus)
ch <- prometheus.MustNewConstMetric(collector.password_expired_status_gauge,prometheus.GaugeValue,vars.Inputs.UserAccountPasswordExpired)
ch <- prometheus.MustNewConstMetric(collector.unknown_authentication_error_status_gauge,prometheus.GaugeValue,vars.Inputs.UserAccountUnknownError)
ch <- prometheus.MustNewConstMetric(collector.credential_validation_status_gauge,prometheus.GaugeValue,vars.Inputs.UserAccountWrongPassword)
}

