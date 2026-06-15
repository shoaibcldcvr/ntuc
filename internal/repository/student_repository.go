package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"ntuc/internal/domain"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(ctx context.Context, student *domain.Student) error {
	query := `
		INSERT INTO students (id, first_name, last_name, email, phone, roll_no, class, gpa, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	student.ID = uuid.New()
	student.CreatedAt = time.Now()
	student.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query, student.ID.String(), student.FirstName, student.LastName, 
		student.Email, student.Phone, student.RollNo, student.Class, student.GPA, 
		student.CreatedAt, student.UpdatedAt)
	return err
}

func (r *StudentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Student, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, roll_no, class, gpa, created_at, updated_at
		FROM students
		WHERE id = ?
	`

	student := &domain.Student{}
	var idStr string
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &student.FirstName, &student.LastName, &student.Email, &student.Phone,
		&student.RollNo, &student.Class, &student.GPA, &student.CreatedAt, &student.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	student.ID, _ = uuid.Parse(idStr)
	return student, nil
}

func (r *StudentRepository) GetByEmail(ctx context.Context, email string) (*domain.Student, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, roll_no, class, gpa, created_at, updated_at
		FROM students
		WHERE email = ?
	`

	student := &domain.Student{}
	var idStr string
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&idStr, &student.FirstName, &student.LastName, &student.Email, &student.Phone,
		&student.RollNo, &student.Class, &student.GPA, &student.CreatedAt, &student.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	student.ID, _ = uuid.Parse(idStr)
	return student, nil
}

func (r *StudentRepository) GetByRollNo(ctx context.Context, rollNo string) (*domain.Student, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, roll_no, class, gpa, created_at, updated_at
		FROM students
		WHERE roll_no = ?
	`

	student := &domain.Student{}
	var idStr string
	err := r.db.QueryRowContext(ctx, query, rollNo).Scan(
		&idStr, &student.FirstName, &student.LastName, &student.Email, &student.Phone,
		&student.RollNo, &student.Class, &student.GPA, &student.CreatedAt, &student.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	student.ID, _ = uuid.Parse(idStr)
	return student, nil
}

func (r *StudentRepository) ListAll(ctx context.Context) ([]*domain.Student, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, roll_no, class, gpa, created_at, updated_at
		FROM students
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*domain.Student
	for rows.Next() {
		student := &domain.Student{}
		var idStr string
		err := rows.Scan(
			&idStr, &student.FirstName, &student.LastName, &student.Email, &student.Phone,
			&student.RollNo, &student.Class, &student.GPA, &student.CreatedAt, &student.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		student.ID, _ = uuid.Parse(idStr)
		students = append(students, student)
	}

	return students, rows.Err()
}

func (r *StudentRepository) ListByClass(ctx context.Context, class string) ([]*domain.Student, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, roll_no, class, gpa, created_at, updated_at
		FROM students
		WHERE class = ?
		ORDER BY roll_no ASC
	`

	rows, err := r.db.QueryContext(ctx, query, class)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*domain.Student
	for rows.Next() {
		student := &domain.Student{}
		var idStr string
		err := rows.Scan(
			&idStr, &student.FirstName, &student.LastName, &student.Email, &student.Phone,
			&student.RollNo, &student.Class, &student.GPA, &student.CreatedAt, &student.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		student.ID, _ = uuid.Parse(idStr)
		students = append(students, student)
	}

	return students, rows.Err()
}

func (r *StudentRepository) Update(ctx context.Context, student *domain.Student) error {
	query := `
		UPDATE students
		SET first_name = ?, last_name = ?, email = ?, phone = ?, class = ?, gpa = ?, updated_at = ?
		WHERE id = ?
	`

	student.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query, student.FirstName, student.LastName, student.Email,
		student.Phone, student.Class, student.GPA, student.UpdatedAt, student.ID.String())
	return err
}

func (r *StudentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM students WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}
