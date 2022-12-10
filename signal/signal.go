package signal

import "github.com/pion/webrtc/v3"

// Join message
type Join struct {
	SID   string                    `json:"sid"`
	UID   string                    `json:"uid"`
	Offer webrtc.SessionDescription `json:"offer"`
	//Config sfu.JoinConfig            `json:"config"`
}

type Trickle struct {
	Target    int                     `json:"target"`
	Candidate webrtc.ICECandidateInit `json:"candidate"`
}

type Negotiation struct {
	Desc webrtc.SessionDescription `json:"desc"`
}

type Signal interface {
	Setonnegotiate(onnegotiate func(webrtc.SessionDescription))
	Setontrickle(ontrickle func(Trickle))
	Join(sid string, uid string, offer webrtc.SessionDescription) webrtc.SessionDescription //offer RTCSessionDescriptionInit  & return RTCSessionDescriptionInit
	Offer(offer webrtc.SessionDescription) webrtc.SessionDescription                        // return RTCSessionDescriptionInit
	Answer()
	//Trickle(t Trickle)
	Trickle(target int, candidate webrtc.ICECandidateInit)
	Close()
}

/*
実装モデル
	onnegotiate?: (jsep: RTCSessionDescriptionInit) => void;
	ontrickle?: (trickle: Trickle) => void;

	join(sid: string, uid: null | string, offer: RTCSessionDescriptionInit): Promise<RTCSessionDescriptionInit>;
	offer(offer: RTCSessionDescriptionInit): Promise<RTCSessionDescriptionInit>;
	answer(answer: RTCSessionDescriptionInit): void;
	trickle(trickle: Trickle): void;
	close(): void;
*/
