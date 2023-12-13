package conf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Conf struct {
	Instances []struct {
		Addr string `json:"addr"`
		Pass string `json:"pass"`
	} `json:"instances"`
}

// configuration parser

func ConfParse(f *os.File) (*Conf, error) {
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read configuration: %v", err)
	}
	conf := &Conf{}
	err = json.Unmarshal(buf, conf)
	if err != nil {
		return nil, fmt.Errorf("configuration unmarshal: %v", err)
	}
	return conf, nil
}
