package application

const defaultBranch = "0001"

type BranchPolicy interface {
	Branch() string
}

type DefaultBranchPolicy struct{}

// NewDefaultBranchPolicy creates a new instance of the DefaultBranchPolicy.
func NewDefaultBranchPolicy() DefaultBranchPolicy {
	return DefaultBranchPolicy{}
}

// Branch returns the default branch code.
func (DefaultBranchPolicy) Branch() string {
	return defaultBranch
}
