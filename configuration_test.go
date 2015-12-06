package webrtc

import (
	"fmt"
	"testing"
)

func init() {
	// SetVerbosity(0)
}

// Ensure the Go "enums" generated in the idiomatic iota const way actually
// match up with actual int values of the underlying native WebRTC Enums.
func checkEnum(t *testing.T, desc string, enum int, expected int) {
	if enum != expected {
		t.Error("Mismatched Enum Value -", desc,
			"\nwas:", enum,
			"\nexpected:", expected)
	}
}

func TestBundlePolicy(t *testing.T) {
	checkEnum(t, "BundlePolicyBalanced",
		int(BundlePolicyBalanced), _cgoBundlePolicyBalanced)
	checkEnum(t, "BundlePolicyMaxCompat",
		int(BundlePolicyMaxCompat), _cgoBundlePolicyMaxCompat)
	checkEnum(t, "BundlePolicyMaxBundle",
		int(BundlePolicyMaxBundle), _cgoBundlePolicyMaxBundle)
}

func TestIceTransportPolicy(t *testing.T) {
	checkEnum(t, "IceTransportPolicyNone",
		int(IceTransportPolicyNone), _cgoIceTransportPolicyNone)
	checkEnum(t, "IceTransportPolicyRelay",
		int(IceTransportPolicyRelay), _cgoIceTransportPolicyRelay)
	checkEnum(t, "IceTransportPolicyAll",
		int(IceTransportPolicyAll), _cgoIceTransportPolicyAll)
}

// TODO: [ED]
/* func TestRtcpMuxPolicy(t *testing.T) {
	checkEnum(t, "RtcpMuxPolicyNegotiate",
		int(RtcpMuxPolicyNegotiate), _cgoRtcpMuxPolicyNegotiate)
	checkEnum(t, "RtcpMuxPolicyRequire",
		int(RtcpMuxPolicyRequire), _cgoRtcpMuxPolicyRequire)
} */

func TestIceServer(t *testing.T) {
	s, err := NewIceServer()
	if nil == err {
		t.Error("NewIceServer should have failed given 0 params",
			s.Urls)
	}
	s, err = NewIceServer("")
	if nil == err {
		t.Error("NewIceServer should have failed given empty urls.")
	}
	s, err = NewIceServer("stun:12345, badurl")
	if nil == err {
		t.Error("NewIceServer should have failed given malformed url.")
	}
	s, err = NewIceServer("stun:12345, stun:ok")
	if nil != err {
		t.Error(err)
	}
	s, err = NewIceServer("stun:a, turn:b")
	if nil != err {
		t.Error(err)
	}
	s, err = NewIceServer("stun:a, turn:b", "alice")
	if nil != err {
		t.Error(err)
	}
	s, err = NewIceServer("stun:a, turn:b", "alice", "secret")
	if nil != err {
		t.Error(err)
	}
	s, err = NewIceServer("stun:a, turn:b", "alice", "secret", "extra")
	if nil != err {
		t.Error("NewIceServer shouldn't fail, only WARN on too many params.")
	}
	fmt.Println(s)
}

func TestNewConfiguration(t *testing.T) {
	config := NewRTCConfiguration()
	if nil == config {
		t.Error("NewRTCConfiguration could not generate basic config.")
	}
	config = NewRTCConfiguration(OptionIceServer("stun:a"))
	if len(config.IceServers) != 1 {
		t.Error("NewRTCConfiguration should have 1 ICE server.")
	}
	config = NewRTCConfiguration(
		OptionIceServer("stun:a"),
		OptionIceServer("stun:b, turn:c"))
	if len(config.IceServers) != 2 {
		t.Error("NewRTCConfiguration should have 2 ICE servers.")
	}
}

func TestIceServerCGO(t *testing.T) {
}
