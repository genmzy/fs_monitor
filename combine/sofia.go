package combine

import (
	"fs_monitor/goesl"
	"fs_monitor/goesl/ev_header"
	"fs_monitor/goesl/ev_name"
	"fs_monitor/parse"
)

type SofiaState int

const (
	UNREGED SofiaState = iota
	TRYING
	REGISTER
	REGED
	UNREGISTER
	FAILED
	FAIL_WAIT
	EXPIRED
	NOREG
	DOWN
	TIMEOUT
)

type Profile struct {
	// filled from `sofia status` summary
	Name string
	Addr string
	// filled from event headers
	Core  string
	Index string
	/*
	 * filled from `sofia status profile xxx` details
	 */
	RtpIp    string
	ExtRtpIp string
	SipIp    string
	ExtSipIp string
	CodecIn  []string
	CodecOut []string
	Media    bool
	LateNeg  bool
	DTMFType string
	// need update summaries
	Reg           int
	CallIn        int
	CallInFailed  int
	CallOut       int
	CallOutFailed int
}

type Gateway struct {
	// filled from `sofia status`
	Name    string
	Profile string
	Uname   string
	Realm   string

	/*
	 * filled from `sofia status gateway xxx` details
	 */
	Passwd   bool
	PingFreq int
	// need update summary
	Up            bool
	UpSec         int
	CallIn        int
	CallOut       int
	CallInFailed  int
	CallOutFailed int
	Reg           bool
}

func ProfileEventFeed(e goesl.Event) *Profile {
	en := e.Name()
	if en < ev_name.CHANNEL_CREATE || en > ev_name.CHANNEL_UUID {
		return nil
	}
	p := &Profile{
		Core: e.CoreUuid(),
	}
	if len(e.Get(ev_header.Sip_Call_ID)) == 0 {
		p.Name = e.Get(ev_header.Caller_Profile_Index)
	} else {
		p.Name = e.Get(ev_header.Sip_Profile_Name)
		raw := e.Get(ev_header.Sip_Profile_Url)
		p.Addr = parse.SofiaUrlParse(raw).Host
	}
	return p
}
