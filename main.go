package main

import (
	"flag"
	"os"
	"syscall"

	"github.com/crowdsecurity/crowdsec/pkg/sqlite"
	daemon "github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var config *blockerConfig

func termHandler(sig os.Signal) error {
	cloudflareCTX, err := newCloudflareContext(config)
	if err != nil {
		log.Fatal(err)
	}

	err = cloudflareCTX.deleteAllRules()
	if err != nil {
		log.Fatalf("error while removing all rules: %s", err)
	}
	return daemon.ErrStop
}

func main() {

	var err error

	configPath := flag.String("c", "", "path to cloudflare-blocker.yaml")
	flag.Parse()

	if configPath == nil || *configPath == "" {
		log.Fatalf("config file required")
	}

	config, err := NewConfig(*configPath)
	if err != nil {
		log.Fatalf("unable to load configuration: %s", err)
	}

	/*Configure logging*/
	if config.LogMode == "file" {
		if config.LogDir == "" {
			config.LogDir = "/var/log/"
		}
		log.SetOutput(&lumberjack.Logger{
			Filename:   config.LogDir + "/cloudflare-blocker.log",
			MaxSize:    500, //megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, //disabled by default
		})
		log.SetFormatter(&log.TextFormatter{TimestampFormat: "02-01-2006 15:04:05", FullTimestamp: true})
	} else if config.LogMode != "stdout" {
		log.Fatalf("log mode '%s' unknown, expecting 'file' or 'stdout'", config.LogMode)
	}

	cloudflareCTX, err := newCloudflareContext(config)
	if err != nil {
		log.Fatal(err)
	}

	dbCTX, err := sqlite.NewSQLite(map[string]string{"db_path": config.DBPath})
	if err != nil {
		log.Fatalf("unable to init sqlite : %v", err)
	}

	if config.Daemon == true {
		go cloudflareCTX.Run(dbCTX, config.updateFrequency)

		daemon.SetSigHandler(termHandler, syscall.SIGTERM)
		//daemon.SetSigHandler(ReloadHandler, syscall.SIGHUP)

		dctx := &daemon.Context{
			PidFileName: config.PidDir + "/cloudflare-blocker.pid",
			PidFilePerm: 0644,
			WorkDir:     "./",
			Umask:       027,
		}

		d, err := dctx.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			return
		}
		defer dctx.Release()

		/*if we are into daemon mode, only process signals*/
		err = daemon.ServeSignals()
		if err != nil {
			log.Errorf("Error: %s", err.Error())
		}
	} else {
		cloudflareCTX.Run(dbCTX, config.updateFrequency)
	}

}
