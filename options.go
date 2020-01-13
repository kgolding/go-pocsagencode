package pocsagencode

type Options struct {
	MaxLen       int
	PreambleBits int
}

type OptionFn func(*Options)

func OptionMaxLen(v int) OptionFn {
	return func(o *Options) {
		o.MaxLen = v
	}
}

func OptionPreambleBits(v int) OptionFn {
	return func(o *Options) {
		o.PreambleBits = v
	}
}
