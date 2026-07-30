package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/ports"
	"errors"

	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &AssignOrdersJob{}

type AssignOrdersJob struct {
	assignOrdersCommandHandler commands.AssignOrdersCommandHandler
}

func NewAssignOrdersJob(
	assignOrdersCommandHandler commands.AssignOrdersCommandHandler) (cron.Job, error) {
	if assignOrdersCommandHandler == nil {
		return nil, errors.New("assignOrdersCommandHandler")
	}
	return &AssignOrdersJob{
		assignOrdersCommandHandler: assignOrdersCommandHandler}, nil
}
func (j *AssignOrdersJob) Run() {
	ctx := context.Background()
	command, err := commands.NewAssignOrdersCommand()
	if err != nil {
		log.Error(err)
	}
	err = j.assignOrdersCommandHandler.Handle(ctx, command)
	if err != nil {
		if errors.Is(err, ports.ErrOrderNotFound) {
			return
		}
		log.Error(err)
	}
}
