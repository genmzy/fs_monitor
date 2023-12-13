package handlers

import (
	"context"
	"fs_monitor/goesl"
	"fs_monitor/goesl/ev_name"
	"time"
)

type Gate struct {
	// warning or error logs
	// hangup => hangup complete stuck checks
	// detected play and get digits but no DTMFs received number
	//
}

var gateEvents []ev_name.EventName = []ev_name.EventName{
	ev_name.DTMF, ev_name.RECORD_START, ev_name.RECORD_STOP, ev_name.API,
	ev_name.BACKGROUND_JOB, ev_name.HEARTBEAT, ev_name.SEND_MESSAGE, ev_name.CUSTOM,
}

var gateSubs []string = []string{
	// box user
	"msg::call_answer",
	"msg::call_status",
}

func (h *Gate) OnConnect(conn *goesl.Connection) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn.Plain(ctx, gateEvents, gateSubs)
	conn.Fslog(ctx, goesl.FslogLevel_WARNING)
}

func (h *Gate) OnDisconnect(conn *goesl.Connection, ev goesl.Event) {
	goesl.Noticef("esl disconnected: %v", ev)
}

func (h *Gate) OnClose(con *goesl.Connection) {
	goesl.Noticef("esl connection closed")
}

func (h *Gate) OnEvent(ctx context.Context, conn *goesl.Connection, e goesl.Event) {
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

func (h *Gate) OnLog(ctx context.Context, con *goesl.Connection, fslog goesl.Fslog) {
}
