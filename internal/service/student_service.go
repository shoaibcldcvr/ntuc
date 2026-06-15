package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"ntuc/internal/domain"
	"ntuc/internal/repository"
	"ntuc/pkg/errors"
	"ntuc/pkg/logger"
)

type StudentService struct {
	repo   *repository.StudentRepository
	logger *logger.Logger
}

func NewStudentService(repo *repository.StudentRepository, logger *logger.Logger) *StudentService {
	return &StudentService{repo: repo, logger: logger}
}

func (s *StudentService) CreateStudent(ctx context.Context, req *domain.CreateStudentRequest) (*domain.StudentResponse, error) {
	// Check if email already exists
	existingStudent, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Errorf("failed to check existing email: %v", err)
		return nil, errors.NewInternalError("failed to create student", err)
	}

	if existingStudent != nil {
		return nil, errors.NewConflictError("student with this email already exists")
	}

	// Check if roll number already exists
	existingStudent, err = s.repo.GetByRollNo(ctx, req.RollNo)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Errorf("failed to check existing roll number: %v", err)
		return nil, errors.NewInternalError("failed to create student", err)
	}

	if existingStudent != nil {
		return nil, errors.NewConflictError("student with this roll number already exists")
	}

	student := &domain.Student{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		RollNo:    req.RollNo,
		Class:     req.Class,
		GPA:       req.GPA,
	}

	if err := s.repo.Create(ctx, student); err != nil {
		s.logger.Errorf("failed to create student: %v", err)
		return nil, errors.NewInternalError("failed to create student", err)
	}

	s.logger.Infof("student created: %s (roll_no: %s)", student.ID, student.RollNo)
	return student.ToResponse(), nil
}

func (s *StudentService) GetStudent(ctx context.Context, id string) (*domain.StudentResponse, error) {
	studentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid student id format")
	}

	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Errorf("failed to get student: %v", err)
		return nil, errors.NewInternalError("failed to get student", err)
	}

	if student == nil {
		return nil, errors.NewNotFoundError("student not found")
	}

	return student.ToResponse(), nil
}

func (s *StudentService) ListStudents(ctx context.Context) ([]*domain.StudentResponse, error) {
	students, err := s.repo.ListAll(ctx)
	if err != nil {
		s.logger.Errorf("failed to list students: %v", err)
		return nil, errors.NewInternalError("failed to list students", err)
	}

	responses := make([]*domain.StudentResponse, len(students))
	for i, student := range students {
		responses[i] = student.ToResponse()
	}

	return responses, nil
}

func (s *StudentService) ListStudentsByClass(ctx context.Context, class string) ([]*domain.StudentResponse, error) {
	students, err := s.repo.ListByClass(ctx, class)
	if err != nil {
		s.logger.Errorf("failed to list students by class: %v", err)
		return nil, errors.NewInternalError("failed to list students", err)
	}

	responses := make([]*domain.StudentResponse, len(students))
	for i, student := range students {
		responses[i] = student.ToResponse()
	}

	return responses, nil
}

func (s *StudentService) UpdateStudent(ctx context.Context, id string, req *domain.UpdateStudentRequest) (*domain.StudentResponse, error) {
	studentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid student id format")
	}

	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Errorf("failed to get student: %v", err)
		return nil, errors.NewInternalError("failed to update student", err)
	}

	if student == nil {
		return nil, errors.NewNotFoundError("student not found")
	}

	// Check if email is already taken by another student
	if req.Email != student.Email {
		existingStudent, err := s.repo.GetByEmail(ctx, req.Email)
		if err != nil && err != sql.ErrNoRows {
			s.logger.Errorf("failed to check email: %v", err)
			return nil, errors.NewInternalError("failed to update student", err)
		}

		if existingStudent != nil {
			return nil, errors.NewConflictError("email already in use")
		}
	}

	student.FirstName = req.FirstName
	student.LastName = req.LastName
	student.Email = req.Email
	student.Phone = req.Phone
	student.Class = req.Class
	student.GPA = req.GPA

	if err := s.repo.Update(ctx, student); err != nil {
		s.logger.Errorf("failed to update student: %v", err)
		return nil, errors.NewInternalError("failed to update student", err)
	}

	s.logger.Infof("student updated: %s", student.ID)
	return student.ToResponse(), nil
}

func (s *StudentService) DeleteStudent(ctx context.Context, id string) error {
	studentID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid student id format")
	}

	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Errorf("failed to get student: %v", err)
		return errors.NewInternalError("failed to delete student", err)
	}

	if student == nil {
		return errors.NewNotFoundError("student not found")
	}

	if err := s.repo.Delete(ctx, studentID); err != nil {
		s.logger.Errorf("failed to delete student: %v", err)
		return errors.NewInternalError("failed to delete student", err)
	}

	s.logger.Infof("student deleted: %s", studentID)
	return nil
}
