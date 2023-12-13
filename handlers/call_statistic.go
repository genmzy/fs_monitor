package handlers

import (
	"context"
	"fs_monitor/combine"
	"fs_monitor/goesl"
)

type CallStatistic interface {
	OnPreLeg(context.Context, goesl.ConnHandler, *combine.PreLeg)
	OnLeg(context.Context, goesl.ConnHandler, *combine.Leg)
	OnLegSummary(context.Context, goesl.ConnHandler, *combine.LegSummary)
	OnLegBridge(context.Context, goesl.ConnHandler, *combine.Bridge)
}
