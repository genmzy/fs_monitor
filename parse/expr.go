package parse

import (
	"fs_monitor/goesl"
	"strings"
	"time"
)

func ExprTrue(expr string) bool {
	expr = strings.ToLower(expr)
	switch expr {
	case "yes":
	case "on":
	case "true":
	case "t":
	case "enabled":
	case "active":
	case "allow":
	case "1":
	default:
		return false
	}
	return true
}

func ExprSec(expr string) int {
	x, err := time.ParseDuration(expr)
	if err != nil {
		goesl.Errorf("invalid expr %s: %v", expr, err)
		return 0
	}
	return int(x / time.Second)
}
