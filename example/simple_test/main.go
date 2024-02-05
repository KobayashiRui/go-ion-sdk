package main

import (
	"go-ion-sdk/client"
	"go-ion-sdk/local_stream"
	"go-ion-sdk/signal"

	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/x264" // This is required to use h264 video encoder
	"github.com/pion/mediadevices/pkg/prop"
)

func main() {

	signal := signal.NewIonSFUJSONRPCSignal()
	signal.Connect("ws://127.0.0.1:7001/ws")

	client := client.NewDefaultClient(signal)

	client.Join("ion", "test uid")

	/*
		//音声圧縮用にOpusの設定を定義
		opusParams, err := opus.NewParams()
		if err != nil {
			panic(err)
		}

		//コーデックを設定する
		codecSelector := mediadevices.NewCodecSelector(
			mediadevices.WithVideoEncoders(&x264Params),
			mediadevices.WithAudioEncoders(&opusParams),
		)
	*/

	/*
		//デバイスリストを取得
		devices := mediadevices.EnumerateDevices()
		for _, device := range devices {
			//fmt.Printf("%v : %v \n", i, device)
			fmt.Println("-----------------")
			fmt.Printf("lavel: %v\n", device.Label)
			fmt.Printf("Device ID: %v\n", device.DeviceID)
			fmt.Printf("Kind: %v\n", device.Kind)
			fmt.Printf("Type: %v\n", device.DeviceType)

			fmt.Println("-----------------")
		}
	*/
	//動画圧縮用のX264の設定を定義
	x264Params, err := x264.NewParams()
	if err != nil {
		panic(err)
	}

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&x264Params),
	)

	x264Params.BitRate = 500_000 // 500kbps

	local := local_stream.NewLocalStream(mediadevices.MediaStreamConstraints{
		Video: func(constraint *mediadevices.MediaTrackConstraints) {
			// Query for ideal resolutions
			constraint.Width = prop.Int(600)
			constraint.Height = prop.Int(400)
		},
		Codec: codecSelector,
	})

	client.Publish(local)

	// Block forever
	select {}
}
