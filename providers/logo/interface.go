package logo

type (
	LogoServicer interface {
		GenerateLogo() (string, error)
	}
)
