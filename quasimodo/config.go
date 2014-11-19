package quasimodo

import (
  "github.com/BurntSushi/toml"
  homedir "github.com/mitchellh/go-homedir"
)

type Config struct {
  ConfigFile           string   `toml:"conf"`
  DataDir              string   `toml:"datadir"`
  WatchIntervalSeconds int64    `toml:"watch_interval_sec"`
  PrepareSeconds       int64    `toml:"prepare_sec"`
  StartOffsetSeconds   int64    `toml:"start_offset_sec"`
  StopOffsetSeconds    int64    `toml:"stop_offset_sec"`
  MaxTasks             int      `toml:"max_tasks"`
  AllowedHosts         []string `toml:"allowed_hosts"`
  Host                 string   `toml:"host"`
  Port                 string   `toml:"port"`
  BasePath             string   `toml:"basepath"`
}

// TODO: Linux only
var DefaultConfig = &Config{
  ConfigFile:           "/etc/qsmd/qsmd.conf",
  DataDir:              "/var/lib/qsmd",
  WatchIntervalSeconds: 10,
  PrepareSeconds:       300,
  StartOffsetSeconds:   5,
  StopOffsetSeconds:    5,
  MaxTasks:             10,
  AllowedHosts:         []string{"127.0.0.1"},
  Host:                 "",
  Port:                 "1831",
  BasePath:             "/_qsmd/",
}

func LoadConfig(file string) (*Config, error) {
  conf, err := LoadConfigFile(file)
  if err != nil {
    return conf, err
  }

  // TODO: Set default to empty settings
  if conf.ConfigFile == "" {
    conf.ConfigFile = file
  }
  if conf.DataDir == "" {
    conf.DataDir = DefaultConfig.DataDir
  }
  if conf.WatchIntervalSeconds < 1 {
    conf.WatchIntervalSeconds = DefaultConfig.WatchIntervalSeconds
  }
  // TODO: int64 default: Zero
  if conf.PrepareSeconds == 0 {
    conf.PrepareSeconds = DefaultConfig.PrepareSeconds
  }
  if len(conf.AllowedHosts) == 0 {
    conf.AllowedHosts = DefaultConfig.AllowedHosts
  }
  if conf.Host == "" {
    conf.Host = DefaultConfig.Host
  }
  if conf.Port == "" {
    conf.Port = DefaultConfig.Port
  }
  if conf.BasePath == "" {
    conf.BasePath = DefaultConfig.BasePath
  }

  dataDir, err := homedir.Expand(conf.DataDir)
  if err != nil {
    panic(err)
  }
  conf.DataDir = dataDir

  return conf, err
}

func LoadConfigFile(file string) (*Config, error) {
  var conf Config
  _, err := toml.DecodeFile(file, &conf)
  if err != nil {
    return &conf, err
  }
  return &conf, nil
}
