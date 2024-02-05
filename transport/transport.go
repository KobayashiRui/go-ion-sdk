package transport

import (
	"fmt"
	"github.com/KobayashiRui/go-ion-sdk/signal"

	"github.com/pion/webrtc/v3"
)

type Role int

const API_CHANNEL = "ion-sfu"

const (
	Role_Pub Role = 0
	Role_Sub Role = 1
)

type Transport struct {
	pc         *webrtc.PeerConnection
	signal     signal.Signal
	candidates []webrtc.ICECandidateInit
}

func NewTransport(role Role, signal signal.Signal, config webrtc.Configuration) *Transport {

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	if role == Role_Pub {
		fmt.Println("Role pub")
		//dataChannel, err := peerConnection.CreateDataChannel(API_CHANNEL, nil)
		_, err := peerConnection.CreateDataChannel(API_CHANNEL, nil)
		if err != nil {
			panic(err)
		}

	} else {
		fmt.Println("Role sub")
	}

	peerConnection.OnNegotiationNeeded(func() {})

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {

		if c != nil {
			fmt.Printf("on ice candidate [%v] \n", role)

			//peerConnection.AddICECandidate(c.ToJSON())
			//t := signal.Trickle{
			//	Target:    0,
			//	Candidate: *c,
			//}

			signal.Trickle(int(role), c.ToJSON())
		}
	})

	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed[%v]: %s\n", role, s.String())
	})

	//peerConnection.OnNegotiationNeeded(func() {})

	return &Transport{
		pc:         peerConnection,
		signal:     signal,
		candidates: make([]webrtc.ICECandidateInit, 0, 0),
	}
}

func (c *Transport) GetPeerConnection() *webrtc.PeerConnection {
	return c.pc
}

func (c *Transport) AddCandidates(candidate webrtc.ICECandidateInit) {
	c.candidates = append(c.candidates, candidate)
}

func (c *Transport) GetCandidates() []webrtc.ICECandidateInit {
	return c.candidates
}
