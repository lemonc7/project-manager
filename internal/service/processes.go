package service

import (
	"context"
	"database/sql"

	"github.com/lemonc7/project-manager/internal/model"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/pkg/util"
)

type ProcessService struct {
	querier repo.Querier
}

func NewProcessService(db *sql.DB) *ProcessService {
	return &ProcessService{
		querier: repo.New(db),
	}
}

func (p *ProcessService) CreateProcess(ctx context.Context, req model.CreateProcessRequest) error {
	err := p.querier.CreateProcess(ctx, repo.CreateProcessParams{
		Name:        req.Name,
		ProjectID:   req.ProjectID,
		EmployeeID:  util.ToNullInt32(req.EmployeeID),
		Description: req.Description,
		Sequence:    req.Sequence,
		Remark:      req.Remark,
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *ProcessService) UpdateProcess(ctx context.Context, req model.UpdateProcessRequest) error {
	var status repo.NullProcessStatus
	if req.Status != nil {
		status.Valid = true
		status.ProcessStatus = repo.ProcessStatus(*req.Status)
	}
	err := p.querier.UpdateProcess(ctx, repo.UpdateProcessParams{
		Name:        util.ToNullString(req.Name),
		Description: util.ToNullString(req.Description),
		Status:      status,
		EmployeeID:  util.ToNullInt32(req.EmployeeID),
		Sequence:    util.ToNullInt32(req.Sequence),
		Remark:      util.ToNullString(req.Remark),
		ID:          req.ID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *ProcessService) DeleteProcess(ctx context.Context, req model.DeleteProcessRequest) error {
	err := p.querier.DeleteProcess(ctx, req.ID)

	if err != nil {
		return err
	}

	return nil
}

func (p *ProcessService) ListProcessesFiltered(ctx context.Context, req model.ListProcessesRequest) ([]repo.ListProcessesFilteredRow, error) {
	var status repo.NullProcessStatus
	if req.Status != nil {
		status.Valid = true
		status.ProcessStatus = repo.ProcessStatus(*req.Status)
	}

	processes, err := p.querier.ListProcessesFiltered(ctx, repo.ListProcessesFilteredParams{
		ProjectID:  util.ToNullInt32(req.ProjectID),
		Status:     status,
		EmployeeID: util.ToNullInt32(req.EmployeeID),
	})

	if err != nil {
		return nil, err
	}

	return processes, nil
}
