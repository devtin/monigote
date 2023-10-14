package models

type Args struct {
	Input      string   `arg:"positional,required"`
	Interfaces []string `arg:"-i,--interfaces"`
	Output     string   `arg:"-o,--output"`
}
