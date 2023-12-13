package main

import "fs_monitor/goesl"

var handlers = make(map[string]*goesl.ConnHandler)

func RegHandlers(h *goesl.ConnHandler, htype string) {
	if _, ok := handlers[htype]; ok {
		goesl.Warnf("duplicate handler type %s", htype)
		return
	}
	handlers[htype] = h
}
