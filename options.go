package dtsfmt

type Option func(p *Prettier)

func WithDebug(debug bool) Option {
	return func(p *Prettier) {
		p.isDebug = debug
	}
}

func WithIntegerCellSize(size int) Option {
	return func(p *Prettier) {
		p.integerCellSize = size
	}
}
