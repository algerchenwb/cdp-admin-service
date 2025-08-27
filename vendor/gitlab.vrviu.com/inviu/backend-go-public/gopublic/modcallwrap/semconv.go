package modcallwrap

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	CallerIDKey    = attribute.Key("modcall.CallerID")
	CalleeIDKey    = attribute.Key("modcall.CalleeID")
	InterfaceIDKey = attribute.Key("modcall.InterfaceID")
	CallerHostKey  = attribute.Key("modcall.CallerHost")
	CalleeHostKey  = attribute.Key("modcall.CalleeHost")
	SessionIDKey   = attribute.Key("modcall.SessionID")
	RequestIDKey   = attribute.Key("modcall.RequestID")
	RetCodeKey     = attribute.Key("modcall.RetCode")
	RetMsgKey      = attribute.Key("modcall.RetMsg")
	TimeCostKey    = attribute.Key("modcall.TimeCost")
	NetSegmentKey  = attribute.Key("modcall.NetSegment")
)
