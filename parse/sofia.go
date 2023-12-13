package parse

import (
	"bufio"
	"bytes"
	"fmt"
	"fs_monitor/logger"
	"net/url"
	"regexp"
	"strings"
)

func SofiaUrlParse(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		logger.Errorf(nil, "parse url %s: %v", raw, err)
		return nil
	}
	if u.User.Username() != "sip" {
		logger.Errorf(nil, "unknown user: %s", u.User.Username())
		return nil
	}
	return u
}

type SofiaEntry struct {
	Name   string
	Type   string // gateway or profile
	Url    string
	Status string
	Extra  string
}

func (se SofiaEntry) String() string {
	return fmt.Sprintf(`{"name":"%s","type":"%s","url":"%s","status":"%s","extra":"%s"}`,
		se.Name, se.Type, se.Url, se.Status, se.Extra)
}

// return sofia entry groups
func SofiaStatus(raw []byte) (entries []*SofiaEntry) {
	s := bufio.NewScanner(bytes.NewBuffer(raw))
	start := false
	for s.Scan() {
		t := s.Text()
		if t[0] == '=' { // ====================
			if start { // end line
				break
			} else { // start line
				start = true
				continue
			}
		}
		if !start {
			continue
		}
		// content line now
		t = strings.TrimSpace(t)
		reg, err := regexp.Compile(`\s+`)
		if err != nil {
			logger.Errorf(nil, "regexp compile %s: %v", t, err)
			continue
		}
		sli := reg.Split(t, -1)
		entry := &SofiaEntry{Name: sli[0], Type: sli[1], Url: sli[2], Status: sli[3]}
		if len(sli) > 4 {
			entry.Extra = sli[4]
		}
		entries = append(entries, entry)
	}
	return
}

// return parameters and value pairs
func SofiaStatusDetails(raw []byte) map[string]string {
	m := make(map[string]string)
	s := bufio.NewScanner(bytes.NewBuffer([]byte(raw)))
	for s.Scan() {
		t := s.Text()
		if t[0] == '=' {
			continue
		}
		sli := strings.Split(s.Text(), "\t")
		if len(sli) != 2 {
			continue
		}
		m[strings.TrimSpace(sli[0])] = strings.TrimSpace(sli[1])
	}
	return m
}
