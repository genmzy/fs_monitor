package handlers

import (
	"context"
	"fs_monitor/combine"
	"fs_monitor/goesl"
	"fs_monitor/goesl/ev_header"
	"fs_monitor/goesl/ev_name"
	"fs_monitor/logger"
	"time"
)

type Base struct {
	// warning or error logs
	// hangup => hangup complete stuck checks
	// detected play and get digits but no DTMFs received number
	Instance *combine.Instance

	Conn  *goesl.Connection
	Ens   []ev_name.EventName
	Subs  []string
	Level goesl.FslogLevel

	PreLegs      map[string]*combine.PreLeg     // uuid => pre-leg
	Legs         map[string]*combine.Leg        // call id => pre-leg
	LegSummaries map[string]*combine.LegSummary // call id => leg-summary
}

var baseEvents = []ev_name.EventName{
	ev_name.CHANNEL_CREATE, ev_name.CHANNEL_ANSWER, ev_name.CHANNEL_HANGUP,
	ev_name.CHANNEL_HANGUP_COMPLETE, ev_name.CHANNEL_PROGRESS, ev_name.CHANNEL_PROGRESS_MEDIA,
	ev_name.CHANNEL_BRIDGE, ev_name.MESSAGE,
	ev_name.DTMF, ev_name.API, ev_name.BACKGROUND_JOB, ev_name.HEARTBEAT,
}

var baseSubs = []string{}

func NewBase() *Base {
	return &Base{
		Ens:      baseEvents,
		Subs:     baseSubs,
		Level:    goesl.FslogLevel_WARNING, // warning logger
		Instance: nil,
	}
}

func (h *Base) OnConnect(conn *goesl.Connection) {
	h.Conn = conn
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	instance := &combine.Instance{
		Conn: conn,
		Addr: conn.Address,
		Pass: conn.Password,
	}
	var err error
	instance.CoreUuid, err = conn.Api(ctx, "global_getvar", "core_uuid")
	if err != nil {
		conn.Close()
		logger.Errorf(instance, "api global_getvar: %v", err)
		return
	}
	// return exit or not
	check := func(err error) bool {
		if err == nil {
			return false
		}
		logger.Errorf(instance, "on connect error: %v", err)
		conn.Close()
		return true
	}
	sofiaFuncs := []func(context.Context) error{
		instance.SofiasFeed,
		instance.SofiaProfilesDetailsFeed,
		instance.SofiaGatewaysDetailFeed,
	}
	for _, sofiaFunc := range sofiaFuncs {
		if check(sofiaFunc(ctx)) {
			return
		}
	}
	conn.Plain(ctx, baseEvents, baseSubs)
	conn.Fslog(ctx, h.Level)
}

func (h *Base) OnDisconnect(conn *goesl.Connection, e goesl.Event) {
	// coreUuid := h.getCoreUuidByAddr(conn.Address)
	// if coreUuid == "" {
	// 	goesl.Errorf("cannot get instance of %s", conn.Address)
	// 	return
	// }
	// h.Instances.Noticef("disconnect %s", coreUuid)
}

/*
 * pre leg feed => channel create
 *
 * leg feed => channel progress media/channel answer -> record start/stop
 *
 * leg summary feed => channel hangup complete
 */

func (h *Base) PreLegCollect(e goesl.Event) {
	pl, err := h.Instance.PreLegFeed(e)
	if err != nil {
		logger.Errorf(h.Instance, "pre leg feed: %v", err)
		return
	}
	h.PreLegs[pl.Uuid] = pl
}

func (h *Base) LegCollect(e goesl.Event) {
	uuid := e.Uuid()
	pl, ok := h.PreLegs[uuid]
	if !ok {
		logger.Errorf(h.Instance, "invalid %v event for uuid: %s", e.Name(), uuid)
		return
	}
	l := pl.LegFeed(e)
	callid := e.Get(ev_header.Sip_Call_ID)
	h.Legs[callid] = l
}

func (h *Base) LegSummaryCollect(e goesl.Event) {
	callid := e.Get(ev_header.Sip_Call_ID)
	l, ok := h.Legs[callid]
	if !ok {
		logger.Errorf(h.Instance, "invalid %v event for call id: %s", e.Name(), callid)
	}
	ls := l.LegSummaryFeed(e)
	h.LegSummaries[l.CallId] = ls
}

func (h *Base) ChannelEventHandle(ctx context.Context, conn *goesl.Connection, e goesl.Event) {
	en := e.Name()
	switch en {
	case ev_name.CHANNEL_CREATE:
		h.PreLegCollect(e)
	case ev_name.CHANNEL_PROGRESS:
		h.LegCollect(e)
	case ev_name.CHANNEL_PROGRESS_MEDIA:
		h.LegCollect(e)
	case ev_name.CHANNEL_ANSWER:
		callid := e.Get(ev_header.Sip_Call_ID)
		l, ok := h.Legs[callid]
		if ok && l.Early { // exist and early media, should not need to collect again
			return
		}
		h.LegCollect(e) // exist but no early media, should collect again for audio informations
	case ev_name.RECORD_START:
		callid := e.Get(ev_header.Sip_Call_ID)
		l, ok := h.Legs[callid]
		if !ok {
			logger.Errorf(h.Instance, "invalid event %v for call id: %s", en, callid)
			return
		}
		l.RecordUpdate(e)
	case ev_name.CHANNEL_HANGUP_COMPLETE:
		h.LegSummaryCollect(e)
	}
}

func (h *Base) OnEvent(ctx context.Context, conn *goesl.Connection, e goesl.Event) {
	if e.Name().IsChannelEvent() {
		h.ChannelEventHandle(ctx, conn, e)
	}
}

func (h *Base) OnFsLog(ctx context.Context, conn *goesl.Connection, fslog goesl.Fslog) {
}

func (h *Base) OnClose(conn *goesl.Connection) {
}
