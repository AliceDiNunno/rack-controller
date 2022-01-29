package postgres

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type projectRepo struct {
	db *gorm.DB
}

type Project struct {
	gorm.Model
	ClusterModel

	ID           uuid.UUID
	DisplayName  string
	Slug         string
	Environments []Environment
	Services     []Service

	EventKey uuid.UUID

	UserID uuid.UUID
}

func (p projectRepo) GetProjectByName(name string) (*domain.Project, *e.Error) {
	var project Project
	err := p.db.Where("display_name = ?", name).First(&project).Error
	if err != nil {
		return nil, e.Wrap(domain.ErrProjectNotFound)
	}

	projectDomain := projectToDomain(project)

	return &projectDomain, nil
}

func (p projectRepo) GetProjectByID(id uuid.UUID) (*domain.Project, *e.Error) {
	var project Project
	err := p.db.Where("id = ?", id).First(&project).Error
	if err != nil {
		return nil, e.Wrap(domain.ErrProjectNotFound)
	}

	projectDomain := projectToDomain(project)

	return &projectDomain, nil
}

func (p projectRepo) GetProjectBySlug(slug string) (*domain.Project, *e.Error) {
	var project Project
	err := p.db.Where("slug = ?", slug).First(&project).Error
	if err != nil {
		return nil, e.Wrap(domain.ErrProjectNotFound)
	}

	projectDomain := projectToDomain(project)

	return &projectDomain, nil
}

func (p projectRepo) GetProjectsByUserId(userId uuid.UUID) ([]domain.Project, *e.Error) {
	var projects []Project
	if err := p.db.Where("user_id = ?", userId).Find(&projects).Error; err != nil {
		return nil, e.Wrap(err)
	}

	return projectsToDomain(projects), nil
}

func (p projectRepo) CreateProject(project domain.Project) *e.Error {
	projectToSave := projectFromDomain(project)

	if err := p.db.Create(&projectToSave).Error; err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (p projectRepo) GetProjectByIDAndKey(id uuid.UUID, key uuid.UUID) (*domain.Project, *e.Error) {
	var project Project
	err := p.db.Where("id = ? AND event_key = ?", id, key).First(&project).Error
	if err != nil {
		return nil, e.Wrap(domain.ErrProjectNotFound)
	}

	projectDomain := projectToDomain(project)

	return &projectDomain, nil
}

func projectsToDomain(project []Project) []domain.Project {
	projectSlice := []domain.Project{}

	for _, p := range project {
		projectSlice = append(projectSlice, projectToDomain(p))
	}

	return projectSlice
}

func projectFromDomain(project domain.Project) Project {
	return Project{
		ID:          project.ID,
		DisplayName: project.DisplayName,
		UserID:      project.UserID,
		Slug:        project.Slug,
		EventKey:    project.EventKey,
	}
}

func projectToDomain(project Project) domain.Project {
	return domain.Project{
		ID:          project.ID,
		DisplayName: project.DisplayName,
		UserID:      project.UserID,
		Slug:        project.Slug,
		EventKey:    project.EventKey,
	}
}

func NewProjectRepo(db *gorm.DB) projectRepo {
	return projectRepo{
		db: db,
	}
}
