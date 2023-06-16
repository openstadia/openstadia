package rtc

import (
	o "github.com/openstadia/openstadia/offer"
	"github.com/pion/mediadevices/pkg/codec"
	"github.com/pion/mediadevices/pkg/codec/openh264"
	"github.com/pion/mediadevices/pkg/codec/vpx"
)

func codecParams(offer o.Offer) codec.VideoEncoderBuilder {
	var params codec.VideoEncoderBuilder

	switch offer.Codec.Type {
	case o.Vp8:
		vp8params, err := vpx.NewVP8Params()
		vp8params.BitRate = offer.Codec.BitRate
		if err != nil {
			panic(err)
		}
		params = &vp8params
	case o.Vp9:
		vp9params, err := vpx.NewVP9Params()
		vp9params.BitRate = offer.Codec.BitRate
		if err != nil {
			panic(err)
		}
		params = &vp9params
	case o.Openh264:
		openh264params, err := openh264.NewParams()
		openh264params.BitRate = offer.Codec.BitRate
		if err != nil {
			panic(err)
		}
		params = &openh264params
	}

	return params
}
