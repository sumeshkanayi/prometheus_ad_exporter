package main

import (
	"gopkg.in/ldap.v2"
	"fmt"
	"flag"
	"os"
	"github.com/y0ssar1an/q"

	"strings"
	"log"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//ldapConnectionError 1 means couldnt connect to ldap
//userAccountActive 1 means user account is disabled


var ldapConnectionError float64
var accountLockOutMetrics=prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "account_lockout_status_gauge",
	Help: "This will return 1 if an account is locked in AD or 0",
})

var accountDisabledMetrics=prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "account_disabled_status_gauge",
	Help: "This will return 1 if an account is disabled in AD or 0",
})

var passwordExpiredMetrics=prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "password_expired_status_gauge",
	Help: "This will return 1 if an account password is expired in AD or 0",
})

var wrongCredentialsMetrics=prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "crdential_validation_status_gauge",
	Help: "This will return 1 if crdentials you provided are wrong  in AD or 0",
})

var unknownLoginErrorMetrics=prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "unknown_authentication_error_status_gauge",
	Help: "This will return 1 if  authentication is failed due to an unknown reason in AD or 0",
})

var ldapConnectionStatusMetrics=prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "ldap_connection_status_gauge",
	Help: "This will return 1 if  LDAP connection is failed or 0",
})



type ldapConnection struct {

	connection ldap.Conn
	error error

}

func init(){

prometheus.MustRegister(accountLockOutMetrics)
prometheus.MustRegister(accountDisabledMetrics)
prometheus.MustRegister(passwordExpiredMetrics)
prometheus.MustRegister(unknownLoginErrorMetrics)
prometheus.MustRegister(wrongCredentialsMetrics)
prometheus.MustRegister(ldapConnectionStatusMetrics)

}

func main(){

	ldapServer:=flag.String("server","ad","LDAP Server Hostname or port")
	ldapServerPort:=flag.String("port","389","LDAP Server port")
	serviceAccountName:=flag.String("user","adminxi","User account to be monitored")
	serviceAccountPassword:=os.Getenv("BIND_PASSWORD")
	flag.Parse()
    connection:=connectToLdap(*ldapServer,*ldapServerPort)
    q.Q(connection)
    if connection.error!=nil{
    	ldapConnectionError=1
    	q.Q(ldapConnectionError)
	}else {
		ldapConnectionError = 0
		log.Printf("LDAP connection successfull , going ghead with binding")
		BindError := bindToLdap(connection.connection, *serviceAccountName, serviceAccountPassword)
		userAccountUnknownError,userAccountPasswordExpired,userAccountDisabled,userAccountLocked,userAccountWrongPassword:=analyzeLdapError(BindError)
        log.Printf("Values are ",userAccountUnknownError,userAccountPasswordExpired,userAccountDisabled,userAccountLocked,userAccountWrongPassword)

		accountDisabledMetrics.Add(userAccountDisabled)
		accountLockOutMetrics.Add(userAccountLocked)
		passwordExpiredMetrics.Add(userAccountPasswordExpired)
		ldapConnectionStatusMetrics.Add(ldapConnectionError)
		unknownLoginErrorMetrics.Add(userAccountUnknownError)
		wrongCredentialsMetrics.Add(userAccountWrongPassword)
		passwordExpiredMetrics.Add(userAccountPasswordExpired)



	}



	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

}

func connectToLdap(ldapServer string,port string)(ldapConnection){
	var ldapConnectionDetails ldapConnection

	connection,err:=ldap.Dial("tcp",ldapServer+":"+port)
	ldapConnectionDetails.error=err

	if err!=nil{
		log.Printf("Error dialling to ldap .Make sure you have provided the corect LDAP server and port details")
		ldapConnectionError=1
		return ldapConnectionDetails




	}



	ldapConnectionDetails.connection=*connection
	return ldapConnectionDetails
}

func bindToLdap(connectionPointer ldap.Conn,userName string,password string) (error)  {
    //q.Q(connectionPointer)
	connectionError:=connectionPointer.Bind(userName,password)


	log.Printf("Trying LDAP binding")
	fmt.Println(connectionError)
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
			log.Printf("Seems either an invalid password for no existing user account used")
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

	return userAccountUnknownError,userAccountPasswordExpired,userAccountDisabled,userAccountLocked,userAccountWrongPassword
}