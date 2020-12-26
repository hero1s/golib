package conf

import (
	"github.com/hero1s/golib/log"
	"testing"
)

type Config struct {
	Name   string
	Number int
}

func TestParseYaml(t *testing.T) {
	c := &Config{}
	if err := AutoParseFile("config.yaml", c); err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Debugf(c.Name)
	log.Debugf("%v", c.Number)
}

func TestParseJson(t *testing.T) {
	c := &Config{}
	if err := AutoParseFile("config.json", c); err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Debugf(c.Name)
	log.Debugf("%v", c.Number)
}

func TestParseToml(t *testing.T) {
	c := &Config{}
	if err := AutoParseFile("config.toml", c); err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Debugf(c.Name)
	log.Debugf("%v", c.Number)
}
