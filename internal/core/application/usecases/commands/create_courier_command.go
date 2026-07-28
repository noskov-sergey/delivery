package commands

import (
	"errors"
	"strings"
)

type CreateCourierCommand struct {
	name  string
	speed int

	valid bool
}

func NewCreateCourierCommand(name string, speed int) (CreateCourierCommand, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return CreateCourierCommand{}, errors.New("name cannot be empty")
	}

	if speed <= 0 {
		return CreateCourierCommand{}, errors.New("speed must be greater than zero")
	}

	return CreateCourierCommand{
		name:  name,
		speed: speed,
		valid: true,
	}, nil
}

func (c *CreateCourierCommand) IsValid() bool { return c.valid }
func (c *CreateCourierCommand) Name() string  { return c.name }
func (c *CreateCourierCommand) Speed() int    { return c.speed }
