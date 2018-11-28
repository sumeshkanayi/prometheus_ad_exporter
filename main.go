package main

import (
	"flag"
	"os"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"ad_exporter/vars"

	"ad_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)












func init(){



}

func main(){

	ldapServer:=flag.String("server","ad","LDAP Server Hostname or port")
	ldapServerPort:=flag.String("port","389","LDAP Server port")
	serviceAccountName:=flag.String("user","gitlab","User account to be monitored")
	serviceAccountPassword:=os.Getenv("BIND_PASSWORD")

	flag.Parse()
	log.Printf("Parsing completed")

	vars.Inputs.UserName=*serviceAccountName
	vars.Inputs.LdapBindPassword=serviceAccountPassword
	vars.Inputs.LdapServer=*ldapServer
	vars.Inputs.LdapServerPort=*ldapServerPort
	vars.Inputs.LdapFullPath=vars.Inputs.LdapServer+":"+vars.Inputs.LdapServerPort
	exporterCollector:=exporter.ADmetricsCollector()
	prometheus.MustRegister(exporterCollector)




	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

}

