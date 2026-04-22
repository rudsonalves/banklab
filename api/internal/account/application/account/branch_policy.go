package application

const defaultBranch = "0001"

type BranchPolicy interface {
	Branch() string
}

type DefaultBranchPolicy struct{}

func NewDefaultBranchPolicy() DefaultBranchPolicy {
	return DefaultBranchPolicy{}
}

func (DefaultBranchPolicy) Branch() string {
	return defaultBranch
}
