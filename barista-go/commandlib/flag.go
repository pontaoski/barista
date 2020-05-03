package commandlib

import (
	flag "github.com/spf13/pflag"
)

type FlagList []Flag

func (fl FlagList) GetFlagSet() *flag.FlagSet {
	var fs flag.FlagSet
	fs.Init("", flag.ContinueOnError)
	for _, flag := range fl {
		flag.Register(&fs)
	}
	return &fs
}

// The flag type.
type Flag interface {
	Long() string
	Short() string
	Usage() string
	Register(*flag.FlagSet)
}

type BoolFlag struct {
	LongFlag  string
	ShortFlag string
	FlagUsage string
	Value     bool
}

func (b BoolFlag) Long() string {
	return b.LongFlag
}

func (b BoolFlag) Short() string {
	return b.ShortFlag
}

func (b BoolFlag) Usage() string {
	return b.FlagUsage
}

func (b BoolFlag) Register(f *flag.FlagSet) {
	f.BoolP(b.LongFlag, b.ShortFlag, b.Value, b.FlagUsage)
}

type StringFlag struct {
	LongFlag  string
	ShortFlag string
	FlagUsage string
	Value     string
}

func (s StringFlag) Long() string {
	return s.LongFlag
}

func (s StringFlag) Short() string {
	return s.ShortFlag
}

func (s StringFlag) Usage() string {
	return s.FlagUsage
}

func (s StringFlag) Register(f *flag.FlagSet) {
	f.StringP(s.LongFlag, s.ShortFlag, s.Value, s.FlagUsage)
}
