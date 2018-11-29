package exporter

import (
	"gopkg.in/ldap.v2"
	"log"
	"github.com/y0ssar1an/q"
	"strings"
	"ad_exporter/vars"
)
var ldapConnectionError float64

//var DefaultTimeout=2*time.Second

type ldapConnection struct {

	connection ldap.Conn
	error error

}

func connectToLdap(ldapServer string,port string)(ldapConnection){
	var ldapConnectionDetails ldapConnection
  ldap.DefaultTimeout=vars.Inputs.LdapConnectionTimeOut
  log.Printf("Dialing LDAP server %s , Binding using %s ",vars.Inputs.LdapFullPath,vars.Inputs.UserName)

	connection,err:=ldap.Dial("tcp",ldapServer+":"+port)
	q.Q(connection)

	ldapConnectionDetails.error=err

	if err!=nil{
		log.Printf("Error dialling to ldap .Make sure you have provided the corect LDAP server and port details")
		ldapConnectionError=1
		vars.Inputs.LdapConnectionStatus=ldapConnectionError
		return ldapConnectionDetails



	}

	log.Printf("Dialing LDAP server successfully completed")
	ldapConnectionDetails.connection=*connection

	return ldapConnectionDetails
}

func bindToLdap(connectionPointer ldap.Conn,userName string,password string) (error)  {
	//q.Q(connectionPointer)
	connectionError:=connectionPointer.Bind(userName,password)


	log.Printf("Binding to LDAP")
	if connectionError!=nil {
		log.Printf("LDAP connection error",connectionError)

	}
	//q.Q(connectionError)
	//q.Q(userName)

	return connectionError

}


func analyzeLdapError(BindError error)(userAccountLocked float64,userAccountDisabled float64,userAccountWrongPassword float64,userAccountPasswordExpired float64, userAccountUnknownError float64){

	if BindError != nil {
		q.Q(BindError.Error())

		ldapError := BindError.Error()
		dataSlice := strings.Split(ldapError, "data ")
		dataSliceCode := strings.Split(dataSlice[1], ",")
		ldapUserAccountErrorCode := dataSliceCode[0]
		q.Q(ldapUserAccountErrorCode)
		switch ldapUserAccountErrorCode {

		case "775":
			log.Printf("Account seems to be Locked")
			userAccountLocked=1

		case "533":
			log.Printf("Account seems to Disabled")
			userAccountDisabled=1
		case "52e":
			log.Printf("Seems either an invalid password or entered a non-existing user account used")
			userAccountWrongPassword=1


		case "773":
			log.Printf("Seems like password for the user account is expired .Get it renewed")
			userAccountPasswordExpired=1
		default:
			log.Printf("Some generic error with the account")
			userAccountUnknownError=1

		}



	}else {
		userAccountUnknownError=0
		userAccountPasswordExpired=0
		userAccountDisabled=0
		userAccountLocked=0
		userAccountWrongPassword=0

	}
	vars.Inputs.UserAccountLocked=userAccountLocked
	vars.Inputs.UserAccountPasswordExpired=userAccountPasswordExpired
	vars.Inputs.UserAccountDisabled=userAccountDisabled
	vars.Inputs.UserAccountUnknownError=userAccountUnknownError
	vars.Inputs.UserAccountWrongPassword=userAccountWrongPassword

	return userAccountUnknownError,userAccountPasswordExpired,userAccountDisabled,userAccountLocked,userAccountWrongPassword
}
