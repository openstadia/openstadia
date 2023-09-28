package loopback

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/gen2brain/malgo"
	"github.com/pion/mediadevices/pkg/driver"
	"github.com/pion/mediadevices/pkg/io/audio"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/mediadevices/pkg/wave"
)

const (
	maxDeviceIDLength = 20
	// TODO: should replace this with a more flexible approach
	sampleRateStep    = 1000
	initialBufferSize = 1024
)

var ctx *malgo.AllocatedContext
var hostEndian binary.ByteOrder
var (
	errUnsupportedFormat = errors.New("the provided audio format is not supported")
)

type microphone struct {
	malgo.DeviceInfo
	chunkChan       chan []byte
	deviceCloseFunc func()

	device    *malgo.Device
	inputProp prop.Media
}

func init() {
	Initialize()
}

// Initialize finds and registers active playback or capture devices. This is part of an experimental API.
func Initialize() {
	backends := []malgo.Backend{
		malgo.BackendWasapi,
	}

	var err error
	ctx, err = malgo.InitContext(backends, malgo.ContextConfig{}, func(message string) {
		// TODO Add debug logging
		//log.Printf("Playback device message: %s\n", message)
	})
	if err != nil {
		panic(err)
	}

	devices, err := ctx.Devices(malgo.Playback)
	if err != nil {
		panic(err)
	}

	for _, device := range devices {
		info, err := ctx.DeviceInfo(malgo.Playback, device.ID, malgo.Shared)
		//fmt.Printf("Info %#v\n", info)
		if err == nil {
			priority := driver.PriorityNormal
			if info.IsDefault > 0 {
				priority = driver.PriorityHigh
			}

			name := device.Name()
			name = strings.Trim(name, "\x00")

			driver.GetManager().Register(newMicrophone(info), driver.Info{
				Label:      device.ID.String(),
				DeviceType: driver.Microphone,
				Priority:   priority,
				Name:       name,
			})
		}
	}

	// Decide which endian
	switch v := *(*uint16)(unsafe.Pointer(&([]byte{0x12, 0x34}[0]))); v {
	case 0x1234:
		hostEndian = binary.BigEndian
	case 0x3412:
		hostEndian = binary.LittleEndian
	default:
		panic(fmt.Sprintf("failed to determine host endianness: %x", v))
	}
}

func newMicrophone(info malgo.DeviceInfo) *microphone {
	return &microphone{
		DeviceInfo: info,
	}
}

func (m *microphone) Open() error {
	fmt.Println("Open")
	m.chunkChan = make(chan []byte, 1)
	return nil
}

func (m *microphone) Close() error {
	fmt.Println("Close")
	if m.deviceCloseFunc != nil {
		m.deviceCloseFunc()
	}
	return nil
}

func (m *microphone) closePlaybackDevice() {
	fmt.Println("closePlaybackDevice")
	if m.device != nil {
		// destory playback device
		m.device.Uninit()
		m.device = nil
	}
}

func (m *microphone) defaultPlaybackDevice(inputProp prop.Media) (*malgo.Device, error) {
	fmt.Printf("%#v\n", inputProp)

	fmt.Println("defaultPlaybackDevice")
	config := malgo.DefaultDeviceConfig(malgo.Loopback)
	var callbacks malgo.DeviceCallbacks

	config.PerformanceProfile = malgo.LowLatency
	config.Capture.Channels = uint32(inputProp.ChannelCount)
	config.SampleRate = uint32(inputProp.SampleRate)
	config.PeriodSizeInMilliseconds = uint32(inputProp.Latency.Milliseconds())
	//FIX: Turn on the microphone with the current device id
	config.Capture.DeviceID = m.ID.Pointer()

	if inputProp.SampleSize == 4 && inputProp.IsFloat {
		config.Capture.Format = malgo.FormatF32
	} else if inputProp.SampleSize == 2 && !inputProp.IsFloat {
		config.Capture.Format = malgo.FormatS16
	} else {
		return nil, errUnsupportedFormat
	}

	onRecvChunk := func(_, chunk []byte, framecount uint32) {
		fmt.Println("onRecvChunk")
		m.chunkChan <- chunk
	}
	callbacks.Data = onRecvChunk

	//onDeviceStop := func() {
	//	go func() {
	//		if !m.muted {
	//			m.Restart()
	//		}
	//	}()
	//}
	//callbacks.Stop = onDeviceStop

	device, err := malgo.InitDevice(ctx.Context, config, callbacks)
	if err != nil {
		return nil, err
	}

	err = device.Start()
	if err != nil {
		return nil, err
	}

	//device.SetAllowPlaybackAutoStreamRouting(true)
	//
	//r := device.ChangeVolume(1.0)
	//log.Println("ChangeVolume result:", r)

	return device, nil
}

func (m *microphone) AudioRecord(inputProp prop.Media) (audio.Reader, error) {
	fmt.Println("AudioRecord")
	m.inputProp = inputProp

	decoder, err := wave.NewDecoder(&wave.RawFormat{
		SampleSize:  inputProp.SampleSize,
		IsFloat:     inputProp.IsFloat,
		Interleaved: inputProp.IsInterleaved,
	})
	if err != nil {
		return nil, err
	}

	device, err := m.defaultPlaybackDevice(inputProp)
	if err != nil {
		log.Println("Init default playback device failed with error:", err)
		return nil, err
	}
	m.device = device

	var closeDeviceOnce sync.Once
	m.deviceCloseFunc = func() {
		closeDeviceOnce.Do(func() {
			m.closePlaybackDevice()

			if m.chunkChan != nil {
				close(m.chunkChan)
				m.chunkChan = nil
			}
		})
	}

	var reader audio.Reader = audio.ReaderFunc(func() (wave.Audio, func(), error) {
		fmt.Println("ReaderFunc")
		chunk, ok := <-m.chunkChan
		if !ok {
			m.deviceCloseFunc()
			return nil, func() {}, io.EOF
		}

		decodedChunk, err := decoder.Decode(hostEndian, chunk, inputProp.ChannelCount)
		// FIXME: the decoder should also fill this information
		switch decodedChunk := decodedChunk.(type) {
		case *wave.Float32Interleaved:
			decodedChunk.Size.SamplingRate = inputProp.SampleRate
		case *wave.Int16Interleaved:
			decodedChunk.Size.SamplingRate = inputProp.SampleRate
		default:
			panic("unsupported format")
		}
		return decodedChunk, func() {}, err
	})

	return reader, nil
}

func (m *microphone) Properties() []prop.Media {
	fmt.Println("Properties")
	var supportedProps []prop.Media
	//logger.Debug("Querying properties")

	var isBigEndian bool
	// miniaudio only uses the host endian
	if hostEndian == binary.BigEndian {
		isBigEndian = true
	}

	for _, format := range m.Formats {
		// FIXME: Currently support 48kHz only. We need to implement a resampler first.
		// for sampleRate := m.MinSampleRate; sampleRate <= m.MaxSampleRate; sampleRate += sampleRateStep {
		sampleRate := 48000
		supportedProp := prop.Media{
			Audio: prop.Audio{
				ChannelCount: int(format.Channels),
				SampleRate:   int(sampleRate),
				IsBigEndian:  isBigEndian,
				// miniaudio only supports interleaved at the moment
				IsInterleaved: true,
				// FIXME: should change this to a less discrete value
				Latency: time.Millisecond * 20,
			},
		}

		switch malgo.FormatType(format.Format) {
		case malgo.FormatF32:
			supportedProp.SampleSize = 4
			supportedProp.IsFloat = true
		case malgo.FormatS16:
			supportedProp.SampleSize = 2
			supportedProp.IsFloat = false
		}

		supportedProps = append(supportedProps, supportedProp)
		// }
	}
	return supportedProps
}
