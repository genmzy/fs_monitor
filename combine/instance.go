package combine

import (
	"context"
	"fmt"
	"fs_monitor/goesl"
	"fs_monitor/goesl/ev_header"
	"fs_monitor/goesl/ev_name"
	"fs_monitor/logger"
	"fs_monitor/parse"
	"strconv"
	"strings"
)

type Instance struct {
	Addr     string
	Pass     string
	CoreUuid string

	Conn     *goesl.Connection
	Profiles map[string]*Profile // key: profile name
	Gateways map[string]*Gateway // key: gateway name
}

// Detail implements logger.DetailLogger.
func (in *Instance) Detail() string {
	return fmt.Sprintf(`{"instance":"%s"}`, in.Addr)
}

// Type implements logger.DetailLogger.
func (*Instance) Type() string {
	return "instance"
}

// test1 profile sip:mod_sofia@172.20.161.167:15060 RUNNING (0)
func (in *Instance) SofiasFeed(ctx context.Context) error {
	raw, err := in.Conn.Api(ctx, "sofia", "status")
	if err != nil {
		return fmt.Errorf("api sofia status: %v", err)
	}
	logger.Noticef(in, "sofia init start: ...")
	entries := parse.SofiaStatus([]byte(raw))

	for _, entry := range entries {
		u := parse.SofiaUrlParse(entry.Url)
		if u == nil {
			continue
		}
		switch entry.Type {
		case "profile":
			if entry.Status != "RUNNING" {
				continue
			}
			p := &Profile{
				Name: entry.Name,
				Addr: u.Host,
			}
			in.Profiles[entry.Name] = p
		case "gateway":
			g := &Gateway{
				Realm: u.Host,
				Uname: u.User.Username(),
			}
			tuple := strings.Split(entry.Name, "::") // profile::gateway_name
			if len(tuple) != 2 {
				continue
			}
			g.Profile = tuple[0]
			g.Name = tuple[1]
			in.Gateways[tuple[1]] = g
		default:
			return fmt.Errorf("unknown sofia stauts type: %v", err)
		}
	}
	return nil
}

func (in *Instance) SofiaProfilesDetailsFeed(ctx context.Context) error {
	for _, p := range in.Profiles {
		raw, err := in.Conn.Api(ctx, "sofia", "status", "profile", p.Name)
		if err != nil {
			return fmt.Errorf("api: sofia status profile %s: %v", p.Name, err)
		}
		paras := parse.SofiaStatusDetails([]byte(raw))
		// p.Name = paras["Name"]
		p.RtpIp = paras["RTP-IP"]
		p.ExtRtpIp = paras["Ext-RTP-IP"]
		p.SipIp = paras["SIP-IP"]
		p.ExtSipIp = paras["Ext-SIP-IP"]
		// p.Addr = parse.SofiaUrlParse(paras["URL"]).Host
		p.CodecIn = strings.Split(paras["CODECS IN"], ",")
		p.CodecOut = strings.Split(paras["CODECS OUT"], ",")
		p.Media = !parse.ExprTrue(paras["NOMEDIA"])
		p.DTMFType = paras["DTMF-MODE"]
		p.LateNeg = parse.ExprTrue(paras["LATE-NEG"])
		// need update summaries
		p.CallIn, _ = strconv.Atoi(paras["CALLS-IN"])
		p.CallOut, _ = strconv.Atoi(paras["CALLS-OUT"])
		p.CallOutFailed, _ = strconv.Atoi(paras["FAILED-CALLS-OUT"])
		p.CallInFailed, _ = strconv.Atoi(paras["FAILED-CALLS-IN"])
		p.Reg, _ = strconv.Atoi(paras["REGISTRATIONS"])
	}
	return nil
}

func (in *Instance) SofiaGatewaysDetailFeed(ctx context.Context) error {
	for _, g := range in.Gateways {
		raw, err := in.Conn.Api(ctx, "sofia", "status", "gateway", g.Name)
		if err != nil {
			return fmt.Errorf("api: sofia status gateway %s: %v", g.Name, err)
		}
		paras := parse.SofiaStatusDetails([]byte(raw))
		// g.Name = paras["Name"]
		// g.Uname = paras["Username"]
		// g.Realm = paras["Realm"]
		// g.Profile = paras["Profile"]

		g.Passwd = parse.ExprTrue(paras["Password"])
		g.PingFreq, _ = strconv.Atoi(paras["PingFreq"])
		// need update summary
		g.UpSec, _ = strconv.Atoi(paras["Uptime"])
		g.Up = strings.ToLower(paras["Status"]) == "up" && g.UpSec > 0
		g.CallIn, _ = strconv.Atoi(paras["CallsIn"])
		g.CallOut, _ = strconv.Atoi(paras["CallsOut"])
		g.CallInFailed, _ = strconv.Atoi(paras["FailedCallsIN"])
		g.CallOutFailed, _ = strconv.Atoi(paras["FailedCallsOUT"])
		g.Reg = paras["State"] == "NOREG"
	}
	return nil
}

func (in *Instance) SofiaProfilesDynamicUpdate(ctx context.Context) error {
	for _, p := range in.Profiles {
		raw, err := in.Conn.Api(ctx, "sofia", "status", p.Name)
		if err != nil {
			return fmt.Errorf("api: sofia status gateway %s: %v", p.Name, err)
		}
		paras := parse.SofiaStatusDetails([]byte(raw))
		// need update summaries
		p.CallIn, _ = strconv.Atoi(paras["CALLS-IN"])
		p.CallOut, _ = strconv.Atoi(paras["CALLS-OUT"])
		p.CallOutFailed, _ = strconv.Atoi(paras["FAILED-CALLS-OUT"])
		p.CallInFailed, _ = strconv.Atoi(paras["FAILED-CALLS-IN"])
		p.Reg, _ = strconv.Atoi(paras["REGISTRATIONS"])
	}
	return nil
}

func (in *Instance) SofiaGatewaysDynamicUpdate(ctx context.Context) error {
	for _, g := range in.Gateways {
		raw, err := in.Conn.Api(ctx, "sofia", "status", "gateway", g.Name)
		if err != nil {
			return fmt.Errorf("api: sofia status gateway %s: %v", g.Name, err)
		}
		paras := parse.SofiaStatusDetails([]byte(raw))
		// need update summary
		g.UpSec, _ = strconv.Atoi(paras["Uptime"])
		g.Up = strings.ToLower(paras["Status"]) == "up" && g.UpSec > 0
		g.CallIn, _ = strconv.Atoi(paras["CallsIn"])
		g.CallOut, _ = strconv.Atoi(paras["CallsOut"])
		g.CallInFailed, _ = strconv.Atoi(paras["FailedCallsIN"])
		g.CallOutFailed, _ = strconv.Atoi(paras["FailedCallsOUT"])
		g.Reg = paras["State"] == "NOREG"
	}
	return nil
}

// event must be channel event and have variable header
func (in *Instance) PreLegFeed(e goesl.Event) (*PreLeg, error) {
	if !e.Name().IsChannelEvent() {
		return nil, fmt.Errorf("not an channel event: %v", e.Name())
	}
	callout := e.CallDirection() == "outbound"
	en := e.Name()
	var pl *PreLeg = nil
	// callin => channel progress media event and later
	// callout => channel create
	if en == ev_name.CHANNEL_OUTGOING || en == ev_name.CHANNEL_CALLSTATE || en == ev_name.CHANNEL_STATE {
		goto no_variables
	}
	pl = &PreLeg{
		Instance: in.Addr,
		Uuid:     e.Uuid(),
		Caller:   e.Caller(),
		Callee:   e.Callee(),
	}
	pl.Callout = callout
	if callout {
		pl.Profile = e.Get(ev_header.Sip_Profile_Name)
		pl.RemoteSip = e.Get(ev_header.Sip_To_Host)
	} else {
		pl.Profile = e.Get(ev_header.Sofia_Profile_Name)
		pl.RemoteSip = fmt.Sprintf("%s:%s", e.Get(ev_header.Sip_Network_Ip), e.Get(ev_header.Sip_Network_Port))
	}
	pl.LocalSip = in.Profiles[pl.Profile].Addr
	pl.PreCodecs = strings.Split(e.Get(ev_header.Rtp_Use_Codec_String), ",")
	pl.IgnoreEarly = !parse.ExprTrue(e.Get(ev_header.Originate_Early_Media)) // callin regard as false always
	return pl, nil

no_variables:
	return nil, fmt.Errorf("event %v have no(or not enough) channel variables", en)
}

func (in *Instance) BridgeFeed(e goesl.Event) (*Bridge, error) {
	if en := e.Name(); en != ev_name.CHANNEL_BRIDGE {
		return nil, fmt.Errorf("invalid event %v", en)
	}
	b := &Bridge{AtInstance: in.Addr, CallTuple: [2]string{}}
	b.AtEpoch, _ = strconv.Atoi(e.Get(ev_header.Event_Date_Timestamp))
	b.CallTuple[0] = e.Get(ev_header.Bridge_A_Unique_ID)
	b.CallTuple[1] = e.Get(ev_header.Bridge_B_Unique_ID)
	return b, nil
}
