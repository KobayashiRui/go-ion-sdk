package client

import (
	"fmt"
	"go-ion-sdk/local_stream"
	"go-ion-sdk/signal"
	. "go-ion-sdk/transport"

	"github.com/pion/webrtc/v3"
)

type Client struct {
	Signal     signal.Signal
	Config     webrtc.Configuration
	Transports map[Role]*Transport
}

func NewDefaultClient(s signal.Signal) *Client {

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	return NewClient(s, config)
}

func NewClient(s signal.Signal, config webrtc.Configuration) *Client {

	_client := &Client{
		Signal:     s,
		Config:     config,
		Transports: make(map[Role]*Transport),
	}

	s.Setontrickle(_client.trickle)
	s.Setonnegotiate(_client.negotiate)

	return _client
}

func (c *Client) Join(sid string, uid string) {
	fmt.Println("JOIN!")

	c.Transports[Role_Pub] = NewTransport(Role_Pub, c.Signal, c.Config)

	c.Transports[Role_Sub] = NewTransport(Role_Sub, c.Signal, c.Config)

	//TODO sub ondatachannel

	offer, err := c.Transports[Role_Pub].GetPeerConnection().CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	if err := c.Transports[Role_Pub].GetPeerConnection().SetLocalDescription(offer); err != nil {
		panic(err)
	}

	answer := c.Signal.Join(sid, uid, offer)
	if err := c.Transports[Role_Pub].GetPeerConnection().SetRemoteDescription(answer); err != nil {
		panic(err)
	}

	fmt.Println("### Candidate ###")
	for i, candidate := range c.Transports[Role_Pub].GetCandidates() {
		fmt.Printf("%v : %v \n", i, candidate)
		c.Transports[Role_Pub].GetPeerConnection().AddICECandidate(candidate)
	}

	//c.Transports[Role_Pub].GetPeerConnection().OnNegotiationNeeded(c.onNegotiationNeeded)
	//c.Transports[Role_Pub].GetPeerConnection().OnNegotiationNeeded(func() {
	//	fmt.Println("######## OnNegotiationNeeded!! #########")
	//})

	fmt.Println("JOIN")
}

// func (c *Client) Publish(stream local_stream.LocalStream, encodingParams []webrtc.RTPEncodingParameters) {
func (c *Client) Publish(stream local_stream.LocalStream) {
	//stream.Publish(c.Transports[Role_Pub], encodingParams)
	fmt.Println("Publish!")
	stream.Publish(c.Transports[Role_Pub])
	c.renegotiate(false)
}

func (c *Client) UnPublish() {

}

func (c *Client) trickle(t signal.Trickle) {
	fmt.Println("TRICKLE ME!!!")

	fmt.Println("Target: ", t.Target)

	//TODO no transports

	if c.Transports[Role(t.Target)].GetPeerConnection().RemoteDescription() != nil {
		if err := c.Transports[Role(t.Target)].GetPeerConnection().AddICECandidate(t.Candidate); err != nil {
			fmt.Println("candidate error")
		}

		fmt.Println("Candidate Add")

	} else {
		fmt.Println("No candidate")
		//c.Transports[Role(t.Target)].candidates = append(c.Transports[Role(t.Target)].candidates, t.Candidate)
		c.Transports[Role(t.Target)].AddCandidates(t.Candidate)
	}
}

func (c *Client) negotiate(s webrtc.SessionDescription) {
	fmt.Println("Negotiate ME!!!")

	c.Transports[Role_Sub].GetPeerConnection().SetRemoteDescription(s)
	for i, candidate := range c.Transports[Role_Sub].GetCandidates() {
		fmt.Printf("Sub %v : %v \n", i, candidate)
		c.Transports[Role_Sub].GetPeerConnection().AddICECandidate(candidate)
	}

	answer, err := c.Transports[Role_Sub].GetPeerConnection().CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	if err := c.Transports[Role_Sub].GetPeerConnection().SetLocalDescription(answer); err != nil {
		panic(err)
	}
}

func (c *Client) onNegotiationNeeded() {
	fmt.Println("########## onNegotiationNeeded ##########")
	//c.renegotiate(false)
}

func (c *Client) renegotiate(iceRestart bool) {
	fmt.Println("########## onNegotiation ##########")
	//offer, err := c.Transports[Role_Pub].GetPeerConnection().CreateOffer(&webrtc.OfferOptions{ICERestart: iceRestart})
	offer, err := c.Transports[Role_Pub].GetPeerConnection().CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	//if err := c.Transports[Role_Pub].GetPeerConnection().SetLocalDescription(webrtc.SessionDescription{}); err != nil {
	if err := c.Transports[Role_Pub].GetPeerConnection().SetLocalDescription(offer); err != nil {
		panic(err)
	}

	answer := c.Signal.Offer(offer)

	if err := c.Transports[Role_Pub].GetPeerConnection().SetRemoteDescription(answer); err != nil {
		panic(err)
	}
	//fmt.Println("###### END ######")

}
