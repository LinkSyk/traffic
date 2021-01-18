package internal

type LoadBlanceAlg interface {
	// LoadBlance()
}

type RoundRoBinAlg struct {
}

func NewRoundRoBinAlg() *RoundRoBinAlg {
	return &RoundRoBinAlg{}
}

type IPHashAlg struct {
}

func NewIPHashAlg() *IPHashAlg {
	return &IPHashAlg{}
}

type AtLeastAlg struct {
}

func NewAtLeastAlg() *AtLeastAlg {
	return &AtLeastAlg{}
}

func ChooseLoadBlance(alg LBAlg) LoadBlanceAlg {
	switch alg {
	case LBRoundRoBin:
		return NewRoundRoBinAlg()
	case LBIPHash:
		return NewIPHashAlg()
	case LBAtLeast:
		return NewAtLeastAlg()
	default:
		return NewRoundRoBinAlg()
	}
}
