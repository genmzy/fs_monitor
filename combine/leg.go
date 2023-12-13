package combine

import (
	"fmt"
	"fs_monitor/goesl"
	"fs_monitor/goesl/ev_header"
	"fs_monitor/goesl/ev_name"
	"strconv"
	"strings"
)

type PreLeg struct {
	Instance string
	// unique id
	Uuid string
	// call messages
	Profile     string
	LocalSip    string
	RemoteSip   string
	Caller      string
	Callee      string
	PreCodecs   []string
	Callout     bool
	IgnoreEarly bool
}

// Detail implements logger.DetailLogger.
func (l *PreLeg) Detail() string {
	return fmt.Sprintf(`{"instance":"%s","uuid":"%s"}`, l.Instance, l.Uuid)
}

// Detail implements logger.DetailLogger.
func (*PreLeg) Type() string {
	return "pre-leg"
}

type Leg struct {
	Pre *PreLeg

	CallId     string
	LocalRtp   string
	RemoteRtp  string
	Codec      string
	DtmfType   string
	RecordPath string
	Record     bool
	Early      bool
}

// Detail implements logger.DetailLogger.
func (l *Leg) Detail() string {
	return fmt.Sprintf(`{"instance":"%s","call_id":"%s","uuid":"%s"}`, l.Pre.Instance, l.CallId, l.Pre.Uuid)
}

// Detail implements logger.DetailLogger.
func (*Leg) Type() string {
	return "leg"
}

type LegSummary struct {
	Leg *Leg

	CreateEpoch        int
	ProgressEpoch      int
	ProgressMediaEpoch int
	AnswerEpoch        int
	BridgeEpoch        int
	HangupEpoch        int

	InPackets           int
	InJitterPackets     int
	InJitterLossRate    float64
	InQualityPercentage float64
	OutPackets          int
}

// Detail implements logger.DetailLogger.
func (l *LegSummary) Detail() string {
	return fmt.Sprintf(`{"instance":"%s","call_id":"%s","uuid":"%s"}`, l.Leg.Pre.Instance, l.Leg.CallId, l.Leg.Pre.Uuid)
}

// Detail implements logger.DetailLogger.
func (*LegSummary) Type() string {
	return "leg-summary"
}

func (pl *PreLeg) LegFeed(e goesl.Event) (l *Leg) {
	l = &Leg{Pre: pl, Early: false, Record: false}
	en := e.Name()
	if !en.IsChannelEvent() {
		return
	}
	switch en {
	case ev_name.CHANNEL_OUTGOING, ev_name.CHANNEL_STATE,
		ev_name.CHANNEL_CALLSTATE, ev_name.CODEC:
		return
	case ev_name.CHANNEL_CREATE:
		if !pl.Callout {
			l.CallId = e.Get(ev_header.Sip_Call_ID)
		}
		return
	case ev_name.CHANNEL_PROGRESS:
		l.CallId = e.Get(ev_header.Sip_Call_ID)
		l.Early = false
		//  TODO: figure out callin progress feed logic
		return
	case ev_name.CHANNEL_PROGRESS_MEDIA:
		l.Early = true
	case ev_name.RECORD_START, ev_name.RECORD_STOP:
		l.RecordUpdate(e)
	case ev_name.MEDIA_BUG_START, ev_name.MEDIA_BUG_STOP:
		if strings.Contains(e.Get(ev_header.Media_Bug_Function), "record") {
			l.Record = true
			l.RecordPath = e.Get(ev_header.Media_Bug_Target)
		}
	case ev_name.CHANNEL_EXECUTE, ev_name.CHANNEL_EXECUTE_COMPLETE:
		app, data := e.App()
		if strings.Contains(app, "record") {
			l.Record = true
			l.RecordPath = data
		}
	}
	l.CallId = e.Get(ev_header.Sip_Call_ID)
	l.LocalRtp = fmt.Sprintf("%s:%s", e.Get(ev_header.Local_Media_Ip), e.Get(ev_header.Local_Media_Port))
	l.RemoteRtp = fmt.Sprintf("%s:%s", ev_header.Remote_Media_Ip, ev_header.Remote_Media_Port)
	l.Codec = e.Get(ev_header.Rtp_Use_Codec_Name)
	l.DtmfType = e.Get(ev_header.Dtmf_Type)
	return
}

func (l *Leg) RecordUpdate(e goesl.Event) {
	l.Record = true
	l.RecordPath = e.Get(ev_header.Record_File_Path)
}

func (l *Leg) LegSummaryFeed(e goesl.Event) (ls *LegSummary) {
	ls = &LegSummary{Leg: l}
	if en := e.Name(); en != ev_name.CHANNEL_HANGUP_COMPLETE && en != ev_name.CHANNEL_DESTROY {
		return
	}
	ls.CreateEpoch, _ = strconv.Atoi(e.Get(ev_header.Start_UEpoch))
	if ls.Leg.Early {
		ls.ProgressMediaEpoch, _ = strconv.Atoi(e.Get(ev_header.Progress_Media_UEpoch))
	} else {
		ls.ProgressEpoch, _ = strconv.Atoi(e.Get(ev_header.Progress_UEpoch))
	}
	ls.AnswerEpoch, _ = strconv.Atoi(e.Get(ev_header.Answer_UEpoch))
	ls.HangupEpoch, _ = strconv.Atoi(e.Get(ev_header.End_UEpoch))
	ls.InPackets, _ = strconv.Atoi(e.Get(ev_header.Rtp_Audio_In_Packet_Count))
	ls.InJitterPackets, _ = strconv.Atoi(e.Get(ev_header.Rtp_Audio_In_Jitter_Packet_Count))
	ls.InJitterLossRate, _ = strconv.ParseFloat(e.Get(ev_header.Rtp_Audio_In_Jitter_Loss_Rate), 64)
	ls.InQualityPercentage, _ = strconv.ParseFloat(e.Get(ev_header.Rtp_Audio_In_Quality_Percentage), 64)
	ls.OutPackets, _ = strconv.Atoi(e.Get(ev_header.Rtp_Audio_Out_Packet_Count))
	return
}
