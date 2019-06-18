// dbmond
// Copyright (C) 2019 gywndi@gmail.com in kakaoBank
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package common

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
)

// // ConfigStr configure string
// var ConfigStr map[string]string

// // ConfigInt configure int
// var ConfigInt map[string]int

// Cfg config
var Cfg *goconfig.ConfigFile

// Log logger object
var Log Logger

var AbsPath, Base, Port, ExporterPort string
var LogLevel, SnapshotInterval, SnapshotTombSec int

var PromApi, PromRulePath, PromWorkPath, Promtool, RecRawPrefix, RecStatPrefix string
var PromApiTimeout int

var DBHost, DBUser, DBPassword, Database string
var ShowSql int

var Timezone string
var AlarmAPI map[string]string

// LoadConfig load config
func LoadConfig() {

	// ConfigStr = make(map[string]string)
	// ConfigInt = make(map[string]int)

	// parameter
	var config string
	flag.StringVar(&config, "config", "config.ini", "configuration")
	flag.Parse()

	var err error
	Cfg, err = goconfig.LoadConfigFile(config)
	if err != nil {
		panic("Load confguration failed")
	}

	// Load string configure
	AbsPath, err = filepath.Abs(filepath.Dir(os.Args[0]))

	Base = Cfg.MustValue("dbmond", "base", "/dbmond")
	Port = Cfg.MustValue("dbmond", "port", ":3333")
	ExporterPort = Cfg.MustValue("dbmond", "exporter_port", ":9104")
	LogLevel = Cfg.MustInt("dbmond", "log_level", 2)

	SnapshotInterval = Cfg.MustInt("dbmond", "snapshot_interval", 3)
	SnapshotTombSec = Cfg.MustInt("dbmond", "snapshot_tomb_sec", 600)

	PromApi = Cfg.MustValue("dbmond", "prom_api", "http://127.0.0.1:9090/prometheus")
	PromApiTimeout = Cfg.MustInt("dbmond", "prom_api_timeout", 500)
	PromRulePath = Cfg.MustValue("dbmond", "prom_rule_path", "rule")
	PromWorkPath = fmt.Sprintf("%s/%s", PromRulePath, "work")
	Promtool = Cfg.MustValue("dbmond", "promtool", "promtool")
	RecRawPrefix = Cfg.MustValue("dbmond", "rec_raw_prefix", "raw")
	RecStatPrefix = Cfg.MustValue("dbmond", "rec_stat_prefix", "stat")

	DBHost = Cfg.MustValue("database", "host", "127.0.0.1")
	DBUser = Cfg.MustValue("database", "user", "root")
	DBPassword = Cfg.MustValue("database", "pass", "pass")
	Database = Cfg.MustValue("database", "db", "db")
	ShowSql = Cfg.MustInt("database", "show_sql", 0)

	// web-hook
	Timezone = strings.TrimSpace(os.Getenv("WEB_HOOK_TIMEZONE"))
	if Timezone != "" {
		log.Printf(`WEB_HOOK_TIMEZONE : "%s"`, Timezone)
	}

	critAPI := strings.TrimSpace(os.Getenv("WEB_HOOK_CRIT"))
	if critAPI != "" {
		log.Printf(`WEB_HOOK_CRIT : "%s"`, critAPI)
	}
	warnAPI := strings.TrimSpace(os.Getenv("WEB_HOOK_WARN"))
	if warnAPI != "" {
		log.Printf(`WEB_HOOK_WARN : "%s"`, warnAPI)
	}

	AlarmAPI = map[string]string{
		"critical": critAPI,
		"warining": warnAPI,
	}

	// Work path create
	if os.MkdirAll(PromWorkPath, os.ModePerm); err != nil {
		PanicIf(err)
	}
	Log.Info("work path", PromWorkPath, "ok")

	// Rule path create
	if os.MkdirAll(PromRulePath, os.ModePerm); err != nil {
		PanicIf(err)
	}
	Log.Info("rule path", PromRulePath, "ok")

	// Log Setting
	Log.SetLogLevel(LogLevel)

	// Prometheus http client setting
	SetPrometheus()

	// Load messages
	LoadMSG()
}
