package handlers

import (
	"context"
	"fs_monitor/goesl"
	"fs_monitor/goesl/ev_name"
	"time"
)

type Gfy struct {
	// warning or error logs
	// hangup => hangup complete stuck checks
	// detected play and get digits but no DTMFs received number
}

var gfyEvents []ev_name.EventName = []ev_name.EventName{
	// channel events
	ev_name.CHANNEL_CREATE, ev_name.CHANNEL_ANSWER, ev_name.CHANNEL_HANGUP,
	ev_name.CHANNEL_HANGUP_COMPLETE, ev_name.CHANNEL_PROGRESS, ev_name.CHANNEL_PROGRESS_MEDIA,
	ev_name.CHANNEL_BRIDGE,
	// DTMF/API
	ev_name.DTMF, ev_name.API, ev_name.BACKGROUND_JOB, ev_name.HEARTBEAT, ev_name.CODEC, ev_name.CUSTOM,
}

var gfySubs []string = []string{
	"msg::queue_in",
	"msg::queue_out",
	"msg::audio_upload_complete",
}

func (h *Gfy) OnConnect(conn *goesl.Connection) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn.Plain(ctx, gfyEvents, gfySubs)
	conn.Fslog(ctx, goesl.FslogLevel_WARNING)
}

func (h *Gfy) OnDisconnect(conn *goesl.Connection, ev goesl.Event) {
	goesl.Noticef("esl disconnected: %v", ev)
}

func (h *Gfy) OnClose(con *goesl.Connection) {
	goesl.Noticef("esl connection closed")
}

func (h *Gfy) OnEvent(ctx context.Context, conn *goesl.Connection, e goesl.Event) {
	en := e.Name()
	app, appData := e.App()
	goesl.Debugf("%s - event %s %s %s\n", e.Uuid(), en, app, appData)
	goesl.Debugf("fire time: %s\n", e.FireTime().StdTime().Format("2006-01-02 15:04:05"))
	switch en {
	case ev_name.BACKGROUND_JOB:
	case ev_name.CHANNEL_ANSWER:
	case ev_name.CHANNEL_HANGUP:
	}
}

func (h *Gfy) OnLog(ctx context.Context, con *goesl.Connection, fslog goesl.Fslog) {
}
