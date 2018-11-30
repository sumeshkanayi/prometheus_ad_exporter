package main

import (
	"flag"
	"os"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sumeshkanayi/prometheus_ad_exporter/vars"

	"github.com/sumeshkanayi/prometheus_ad_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)








func init(){



}

func main(){

	ldapServer:=flag.String("server","ad","LDAP Server Hostname or port")
	ldapServerPort:=flag.String("port","389","LDAP Server port")
	serviceAccountName:=flag.String("user","user@domainFQDN","User account to be monitored")
	exporterPort:=flag.String("exporterPort","2134","Port on which the exporter should listen")
	serviceAccountPassword:=os.Getenv("BIND_PASSWORD")
	ldapBindTimeOut:=flag.Duration("timeout",5*time.Second,"How long should the exporter wait before exporter times out")

	flag.Parse()
	log.Printf("Parsing input variables completed")
	log.Printf("Waiting for prometheus to start scrapping")
	if serviceAccountPassword==""{
		log.Printf("Empty password detected .Exiting for now.Please set the BIND_PASSWORD environment variable")
		return
	}

	vars.Inputs.UserName=*serviceAccountName
	vars.Inputs.LdapBindPassword=serviceAccountPassword
	vars.Inputs.LdapServer=*ldapServer
	vars.Inputs.LdapServerPort=*ldapServerPort
	vars.Inputs.LdapFullPath=vars.Inputs.LdapServer+":"+vars.Inputs.LdapServerPort
	vars.Inputs.LdapConnectionTimeOut=*ldapBindTimeOut
	exporterCollector:=exporter.ADmetricsCollector()
	prometheus.MustRegister(exporterCollector)




	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+*exporterPort, nil)

}

