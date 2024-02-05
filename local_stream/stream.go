package local_stream

import (
	//"native-test/ion_sfu_client"
	"fmt"
	"github.com/KobayashiRui/go-ion-sdk/transport"

	"github.com/pion/mediadevices"
	"github.com/pion/webrtc/v3"

	_ "github.com/pion/mediadevices/pkg/driver/camera"     // This is required to register camera adapter
	_ "github.com/pion/mediadevices/pkg/driver/microphone" // This is required to register microphone adapter
)

type LocalStream interface {
	//Publish(transport *transport.Transport, encodingParams []webrtc.RTPEncodingParameters)
	Publish(transport *transport.Transport)
	UnPublish()
}

type localStream struct {
	mediaStream    mediadevices.MediaStream
	pc             *webrtc.PeerConnection
	encodingParams []webrtc.RTPEncodingParameters
	constraints    mediadevices.MediaStreamConstraints
}

func NewLocalStream(constraints mediadevices.MediaStreamConstraints) *localStream {
	stream, err := mediadevices.GetUserMedia(constraints)

	if err != nil {
		panic(err)
	}

	return &localStream{
		constraints: constraints,
		mediaStream: stream,
	}
}

func (l *localStream) publishTrack(track mediadevices.Track) {
	//pcが存在していることをチェック
	if l.pc != nil {
		init := webrtc.RtpTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionSendonly,
		}

		//transceiver, err := l.pc.AddTransceiverFromTrack(track, init)
		fmt.Println("@@@@@ Add track @@@@@")
		transceiver, err := l.pc.AddTransceiverFromTrack(track, init)
		if err != nil {
			panic(err)
		}
		if transceiver.Receiver() != nil {
			fmt.Printf("Transceiver shouldn't have a receiver")
		}

		if transceiver.Sender() == nil {
			fmt.Printf("Transceiver should have a sender")
		}

		if len(l.pc.GetTransceivers()) != 1 {
			fmt.Printf("PeerConnection should have one transceiver but has %d", len(l.pc.GetTransceivers()))
		}

		if len(l.pc.GetSenders()) != 1 {
			fmt.Printf("PeerConnection should have one sender but has %d", len(l.pc.GetSenders()))
		}

	}
}

// func (l *localStream) Publish(transport *transport.Transport, encodingParams []webrtc.RTPEncodingParameters) {
func (l *localStream) Publish(transport *transport.Transport) {
	l.pc = transport.GetPeerConnection()
	//l.encodingParams = encodingParams
	fmt.Println("Get Tracks")
	for _, track := range l.mediaStream.GetTracks() {
		fmt.Printf("%v : %v\n", track.StreamID(), track.Kind())
		l.publishTrack(track)
	}
}

func (l *localStream) UnPublish() {
	fmt.Println("un publish")
}
