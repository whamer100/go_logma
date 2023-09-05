package logma

import "testing"

func TestLogging(t *testing.T) {
	Init()
	SetVerbosity(false)
	Info("This is an info.")
	Warn("This is a warn. I should be Yellow")
	Fatal("AAAA !PANIC IM RED")
	Debug("I should not show up. If I do something FUCKED has happened.")
	DebugRaw(" - I should also not show up.")
}

func TestDebug(t *testing.T) {
	Init()
	Info("Testing Debug stuff.")
	SetVerbosity(true)
	Debug("Holy shit %s, I'm visible!", "Lois")
	DebugRaw(" - Isn't that neat?")
}
