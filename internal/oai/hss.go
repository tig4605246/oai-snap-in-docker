package oai

import (
	"errors"
	"oai-snap-in-docker/internal/pkg/util"
	"os"
	"strings"
)

// StartHss : Start HSS as a daemon
func startHss(OaiObj Oai) error {
	// Get working path, Hostname
	hssConf := OaiObj.Conf.ConfigurationPathofCN + "hss.conf"
	hssFdConf := OaiObj.Conf.ConfigurationPathofCN + "hss_fd.conf"
	hssBin := OaiObj.Conf.SnapBinaryPath + "oai-cn.hss"
	hostname, _ := os.Hostname()
	// Strat configuring oai-hss
	OaiObj.Logger.Print("Configure hss.conf")
	//Replace MySQL address
	retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", "s/127.0.0.1/"+OaiObj.Conf.MysqlDomainName+"/g", hssConf)
	if retStatus.Exit != 0 {
		return errors.New("Set mysql IP in " + hssConf + " failed")
	}
	// Replace Identity
	OaiObj.Logger.Print("Configure hss_fd.conf")
	identity := hostname + ".openair4G.eur" // use the Hostname we got before
	syntax := "s/ubuntu.openair4G.eur/" + identity + "/g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", syntax, hssFdConf)
	if retStatus.Exit != 0 {
		return errors.New("Set realm in " + hssFdConf + " failed")
	}
	// Init hss
	if OaiObj.Conf.Test == false {
		OaiObj.Logger.Print("Init hss")
		retStatus = util.RunCmd(OaiObj.Logger, hssBin+"-init")
		for {
			fail := false
			for i := 0; i < len(retStatus.Stderr); i++ {
				if strings.Contains(retStatus.Stderr[i], "ERROR") {
					OaiObj.Logger.Println("Init error, re-run again")
					fail = true
				}
			}
			if fail {
				retStatus = util.RunCmd(OaiObj.Logger, hssBin+"-init")
			} else {
				break
			}
		}
		// oai-cn.hss-start
		OaiObj.Logger.Print("start hss as daemon")
		util.RunCmd(OaiObj.Logger, hssBin+"-start")
	}
	return nil
}
