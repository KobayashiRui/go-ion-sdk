package signal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/sourcegraph/jsonrpc2"
	websocketjsonrpc2 "github.com/sourcegraph/jsonrpc2/websocket"
)

// ion sfu json rpc signalとの接続用
type IonSFUJSONRPCSignal struct {
	jrpc        *jsonrpc2.Conn
	Ontrickle   func(Trickle)
	Onnegotiate func(webrtc.SessionDescription)
}

func NewIonSFUJSONRPCSignal() *IonSFUJSONRPCSignal {
	return &IonSFUJSONRPCSignal{
		jrpc:        nil,
		Ontrickle:   nil,
		Onnegotiate: nil,
	}
}

// websocketへの接続
func (sig *IonSFUJSONRPCSignal) Connect(uri string) {
	ctx := context.Background()
	fmt.Printf("signal uri:%v\n", uri)

	//TODO Error対処
	c, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
	}

	stream := websocketjsonrpc2.NewObjectStream(c)
	sig.jrpc = jsonrpc2.NewConn(ctx, stream, sig)

}

func (sig *IonSFUJSONRPCSignal) SetConnect(jrcp *jsonrpc2.Conn) {
	sig.jrpc = jrcp
}

func (sig *IonSFUJSONRPCSignal) Join(sid string, uid string, offer webrtc.SessionDescription) webrtc.SessionDescription {
	join := Join{
		SID:   sid,
		UID:   uid,
		Offer: offer,
	}

	ctx := context.Background()

	sdp := webrtc.SessionDescription{}
	if err := sig.jrpc.Call(ctx, "join", join, &sdp); err != nil {
		fmt.Println(err)
	}

	return sdp

}

func (sig *IonSFUJSONRPCSignal) Setonnegotiate(onnegotiate func(webrtc.SessionDescription)) {
	sig.Onnegotiate = onnegotiate
}

func (sig *IonSFUJSONRPCSignal) Setontrickle(ontrickle func(Trickle)) {
	sig.Ontrickle = ontrickle
}

func (sig *IonSFUJSONRPCSignal) Offer(offer webrtc.SessionDescription) webrtc.SessionDescription {
	negotiation := Negotiation{
		Desc: offer,
	}

	ctx := context.Background()

	sdp := webrtc.SessionDescription{}
	if err := sig.jrpc.Call(ctx, "offer", negotiation, &sdp); err != nil {
		fmt.Println(err)
	}

	return sdp
}

func (sig *IonSFUJSONRPCSignal) Answer() {}

// func (sig *IonSFUJSONRPCSignal) Trickle(t Trickle) {
func (sig *IonSFUJSONRPCSignal) Trickle(target int, candidate webrtc.ICECandidateInit) {
	//Notify signal

	fmt.Println("$$$$ Trickle $$$$")
	ctx := context.Background()

	//err := sig.jrpc.Notify(ctx, "trickle", t)
	err := sig.jrpc.Notify(ctx, "trickle", Trickle{
		Target:    target,
		Candidate: candidate,
	})
	if err != nil {
		fmt.Println("ERROR")
	}

}
func (sig *IonSFUJSONRPCSignal) Close() {}

func (sig *IonSFUJSONRPCSignal) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	switch req.Method {
	case "trickle":
		fmt.Printf("trickle! \n")
		var trickle Trickle
		err := json.Unmarshal(*req.Params, &trickle)
		if err != nil {
			fmt.Println("ERROR trickle")
		}

		if sig.Ontrickle != nil {
			sig.Ontrickle(trickle)
		}
		// on trickle?
	case "offer":
		fmt.Printf("offer \n")
		var s webrtc.SessionDescription
		err := json.Unmarshal(*req.Params, &s)
		if err != nil {
			fmt.Println("ERROR offer")
		}
		// onnegotiate
		if sig.Onnegotiate != nil {
			sig.Onnegotiate(s)
		}
	default:
		fmt.Printf("Default : %v \n", req.Method)

	}
}
