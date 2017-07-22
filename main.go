package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kardianos/service"
)

var (
	logger service.Logger
)

type program struct {
	exitCh chan bool

	domain   string
	record   string
	recordID string
	token    string
	interval int
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	p.process()
	go p.watch()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	logger.Info("Exiting...")
	p.exitCh <- true
	return nil
}

func (p *program) watch() {
	for {
		select {
		case <-time.After(time.Duration(p.interval) * time.Second):
			go p.process()
		case <-p.exitCh:
			logger.Info("Exited")
			return
		}
	}
}

func (p *program) process() {
	var (
		currentIP string
		recordIP  string
		err       error
	)
	if currentIP, err = getCurrentIP(); err != nil {
		logger.Errorf("Erroring in get current public ip: %v", err)
		return
	}
	if recordIP, err = getRecordIP(p.record, p.domain); err != nil {
		logger.Errorf("Erroring in get current domain record ip: %v", err)
		return
	}

	logger.Infof("Getting current public ip: %s, record ip: %s", currentIP, recordIP)

	if currentIP == recordIP {
		logger.Info("Current public ip same as the record ip, don't need update")
		return
	}

	logger.Infof("Current public ip updated, now update the record ip with: %s", currentIP)
	err = updateRecordIP(p.token, p.record, p.domain, p.recordID, currentIP)
	if err != nil {
		logger.Errorf("Update record ip failed: %v", err)
		return
	}
	logger.Infof("Update record ip successful")
}

func main() {

	var (
		install   bool
		uninstall bool
	)

	svcConfig := &service.Config{
		Name:        "ddnsclient",
		DisplayName: "DDNS client for dnspod",
		Description: "DDNS client service for watch ip changed and update dnspod managed domain A record.",
	}

	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "-install") ||
			strings.HasPrefix(arg, "-uninstall") ||
			i == 0 {
			continue
		}
		svcConfig.Arguments = append(svcConfig.Arguments, arg)
	}

	prg := &program{
		exitCh: make(chan bool),
	}

	flag.StringVar(&prg.token, "token", "id,token", "dnspod token for api auth")
	flag.IntVar(&prg.interval, "interval", 60, "interval in second for check public ip changed")
	flag.StringVar(&prg.domain, "domain", "example.com", "ddns domain")
	flag.StringVar(&prg.record, "record", "test.ddns", "ddns domain record to update")
	flag.BoolVar(&install, "install", false, "install as system service")
	flag.BoolVar(&uninstall, "uninstall", false, "remove it from system service")

	flag.Parse()

	svc, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if install {
		if err = svc.Install(); err != nil {
			log.Fatalf("Erroring install as service: %v", err)
		}
		return
	}

	if uninstall {
		if err = svc.Uninstall(); err != nil {
			log.Fatalf("Erroring uninstall ddns service: %v", err)
		}
		return
	}

	if logger, err = svc.Logger(nil); err != nil {
		log.Fatal(err)
	}

	recordID, err := getRecordID(prg.token, prg.record, prg.domain)
	if err != nil {
		logger.Errorf("Erroring load config: %v", err)
		return
	}
	logger.Infof("Getting record id: %s", recordID)
	prg.recordID = recordID

	err = svc.Run()
	if err != nil {
		logger.Errorf("Erroring run svc: %v", err)
	}
}
