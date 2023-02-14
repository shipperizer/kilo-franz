package sasl

type SASLType = int

const (
	PlainSASL SASLType = iota
	ScramSASL
)
