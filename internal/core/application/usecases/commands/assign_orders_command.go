package commands

type AssignOrdersCommand struct {
	valid bool
}

func NewAssignOrdersCommand() (AssignOrdersCommand, error) {
	return AssignOrdersCommand{
		valid: true,
	}, nil
}

func (a *AssignOrdersCommand) IsValid() bool { return a.valid }
