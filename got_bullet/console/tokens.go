package console

type Color interface {
	Col() string
}
type Token interface {
	Name() string
}

type TokenTextTertiary struct{}

func (t TokenTextTertiary) Name() string {
	return "tertiary"
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

type TokenAlias struct{}

func (t TokenAlias) Name() string {
	return "alias"
}

type TokenNote struct{}

func (t TokenNote) Name() string {
	return "note"
}

type TokenGroup struct{}

func (t TokenGroup) Name() string {
	return "group"
}

type TokenAlert struct{}

func (t TokenAlert) Name() string {
	return "alert"
}
