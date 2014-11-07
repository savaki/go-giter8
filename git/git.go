package git

type Git struct {
	Git     string
	Target  string
	Verbose bool
}

func New(git, target string) *Git {
	return &Git{
		Git:    git,
		Target: target,
	}
}
