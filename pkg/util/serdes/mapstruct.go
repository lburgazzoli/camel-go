package serdes

import "github.com/mitchellh/mapstructure"

func DecodeStruct(input interface{}, result interface{}) error {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           result,

		// custom hooks
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		),
	})

	if err != nil {
		return err
	}

	if err := dec.Decode(input); err != nil {
		return err
	}

	return nil
}
