package handlers

import (
	"context"
	"fs_monitor/combine"
	"fs_monitor/goesl"
)

type Trace struct {
	bases []Base
}

func (t *Trace) OnPreLeg(ctx context.Context, h goesl.ConnHandler, pl *combine.PreLeg) {
}

func (t *Trace) OnLeg(ctx context.Context, h goesl.ConnHandler, l *combine.Leg) {
}

func (t *Trace) OnLegSummary(ctx context.Context, h goesl.ConnHandler, ls *combine.LegSummary) {
}

func (t *Trace) OnLegBridge(ctx context.Context, h goesl.ConnHandler, b *combine.Bridge) {
}
