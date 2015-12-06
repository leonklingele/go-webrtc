/*
Package webrtc is a golang wrapper on native code WebRTC.

To provide an easier experience for users of this package, there are differences
inherent in the interface written here and the original native code WebRTC. This
allows users to use WebRTC in a more idiomatic golang way. For example, callback
mechanism has a layer of indirection that allows goroutines instead.

The interface here is based mostly on: w3c.github.io/webrtc-pc

There is also a complication in building the dependent static library for this
to work. Furthermore it is possible that this will break on future versions
of libwebrtc, because the interface with the native code is be fragile.

TODO(keroserene): More package documentation, and more documentation in general.
*/
package webrtc

/*
#cgo CXXFLAGS: -std=c++0x
#cgo linux,amd64 pkg-config: webrtc-linux-amd64.pc
#cgo darwin,amd64 pkg-config: webrtc-darwin-amd64.pc
#include "cpeerconnection.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"github.com/keroserene/go-webrtc/datachannel"
	"unsafe"
	// "io"
	"io/ioutil"
	"log"
	"os"
)

var (
	INFO  log.Logger
	WARN  log.Logger
	ERROR log.Logger
	TRACE log.Logger
)

// Logging verbosity level, from 0 (nothing) upwards.
func SetVerbosity(level int) {
	// handle io.Writer
	infoOut := ioutil.Discard
	warnOut := ioutil.Discard
	errOut := ioutil.Discard
	traceOut := ioutil.Discard

	// TODO: Better logging levels
	if level > 0 {
		errOut = os.Stdout
	}
	if level > 1 {
		warnOut = os.Stdout
	}
	if level > 2 {
		infoOut = os.Stdout
	}
	if level > 3 {
		traceOut = os.Stdout
	}

	INFO = *log.New(infoOut,
		"INFO: ",
		// log.Ldate|log.Ltime|log.Lshortfile)
		log.Lshortfile)
	WARN = *log.New(warnOut,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	ERROR = *log.New(errOut,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	TRACE = *log.New(traceOut,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func init() {
	SetVerbosity(0)
}

type SDPHeader struct {
	// Keep track of both a pointer to the C++ SessionDescription object,
	// and the serialized string version (which native code generates)
	cgoSdp      C.CGOsdp
	description string
}

type PeerConnection struct {
	localDescription *SDPHeader
	// currentLocalDescription
	// pendingLocalDescription

	remoteDescription *SDPHeader
	// currentRemoteDescription
	// pendingRemoteDescription

	// addIceCandidate func()
	// signalingState  RTCSignalingState
	// iceGatheringState  RTCIceGatheringState
	// iceConnectionState  RTCIceConnectionState
	canTrickleIceCandidates bool
	// getConfiguration
	// setConfiguration
	// close
	OnIceCandidate func()

	// Event handlers:
	// onnegotiationneeded
	// onicecandidate
	// onicecandidateerror
	// onsignalingstatechange
	// onicegatheringstatechange
	// oniceconnectionstatechange

	cgoPeer C.CGOPeer // Native code internals
}

// PeerConnection constructor.
func NewPeerConnection(config *RTCConfiguration) (*PeerConnection, error) {
	pc := new(PeerConnection)
	pc.cgoPeer = C.CGOInitializePeer() // internal CGO Peer.
	if nil == pc.cgoPeer {
		return pc, errors.New("PeerConnection: failed to initialize.")
	}
	cConfig := config.CGO() // Convert for CGO
	if 0 != C.CGOCreatePeerConnection(pc.cgoPeer, &cConfig) {
		return nil, errors.New("PeerConnection: could not create from config.")
	}
	INFO.Println("Created PeerConnection: ", pc, pc.cgoPeer)
	return pc, nil
}

// CreateOffer prepares an SDP "offer" message, which should be sent to the target
// peer over a signalling channel.
func (pc *PeerConnection) CreateOffer() (*SDPHeader, error) {
	sdp := C.CGOCreateOffer(pc.cgoPeer)
	if nil == sdp {
		return nil, errors.New("CreateOffer: could not prepare SDP offer.")
	}
	offer := new(SDPHeader)
	offer.cgoSdp = sdp
	offer.description = C.GoString(C.CGOSerializeSDP(sdp))
	return offer, nil
}

func (pc *PeerConnection) SetLocalDescription(sdp *SDPHeader) error {
	r := C.CGOSetLocalDescription(pc.cgoPeer, sdp.cgoSdp)
	if 0 != r {
		return errors.New("SetLocalDescription failed.")
	}
	pc.localDescription = sdp
	return nil
}

// readonly localDescription
func (pc *PeerConnection) LocalDescription() (sdp *SDPHeader) {
	return pc.localDescription
}

func (pc *PeerConnection) SetRemoteDescription(sdp *SDPHeader) error {
	r := C.CGOSetRemoteDescription(pc.cgoPeer, sdp.cgoSdp)
	if 0 != r {
		return errors.New("SetRemoteDescription failed.")
	}
	pc.remoteDescription = sdp
	return nil
}

// readonly remoteDescription
func (pc *PeerConnection) RemoteDescription() (sdp *SDPHeader) {
	return pc.remoteDescription
}

// CreateAnswer prepares an SDP "answer" message, which should be sent in
// response to a peer that has sent an offer, over the signalling channel.
func (pc *PeerConnection) CreateAnswer() (*SDPHeader, error) {
	sdp := C.CGOCreateAnswer(pc.cgoPeer)
	if nil == sdp {
		return nil, errors.New("CreateAnswer failed: could not prepare SDP offer.")
	}
	answer := new(SDPHeader)
	answer.cgoSdp = sdp
	answer.description = C.GoString(C.CGOSerializeSDP(sdp))
	return answer, nil
}

// TODO: Above methods blocks until success or failure occurs. Maybe there should
// actually be a callback version, so the user doesn't have to make their own
// goroutine.

func (pc *PeerConnection) CreateDataChannel(label string, dict datachannel.Init) (
	*datachannel.DataChannel, error) {
	cDC := C.CGOCreateDataChannel(pc.cgoPeer, C.CString(label), unsafe.Pointer(&dict))
	if nil == cDC {
		return nil, errors.New("Failed to CreateDataChannel")
	}
	dc := datachannel.New()
	return dc, nil
}

// Install a handler for receiving ICE Candidates.
// func OnIceCandidate(pc PeerConnection) {
// }

func unused() {
	fmt.Println("nothing yet")
}
