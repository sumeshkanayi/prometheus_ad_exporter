package vars

import "time"

type Exporterattributes struct{

	UserName string
	LdapServer string
	LdapServerPort string
	LdapBindPassword string
	LdapFullPath string
	UserAccountUnknownError float64
	UserAccountPasswordExpired float64
	UserAccountDisabled float64
	UserAccountLocked float64
	UserAccountWrongPassword float64
	LdapConnectionStatus float64
	LdapConnectionTimeOut time.Duration

}
var Inputs Exporterattributes
