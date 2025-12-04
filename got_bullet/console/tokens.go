package console

type Color interface {
	Col() string
}
type Token interface {
	Name() string
}

type TokenSecondary struct{}

func (t TokenSecondary) Name() string {
	return "secondary"
}

type TokenPrimary struct{}

func (t TokenPrimary) Name() string {
	return "primary"
}

type TokenBrand struct{}

func (t TokenBrand) Name() string {
	return "brand"
}

type TokenGid struct{}

func (t TokenGid) Name() string {
	return "gid"
}

type TokenComplete struct{}

func (t TokenComplete) Name() string {
	return "complete"
}
