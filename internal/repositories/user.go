package repositories

import (
	"context"
	"database/sql"
	"gogrpccrud/pkg/api"

	"time"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServise struct {
	api.UnimplementedUserServer
	db *sql.DB
}

func NewUserServise(db *sql.DB) *UserServise {
	return &UserServise{
		db: db,
	}
}

func (s *UserServise) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	var id string

	err := s.db.QueryRowContext(ctx, `
        INSERT INTO users (full_name, age, phone_number, habits, alive, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `, req.FullName, req.Age, req.PhoneNumber, req.Habits, req.Alive, time.Now()).Scan(&id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &api.CreateUserResponse{
		Id:          id,
		FullName:    req.FullName,
		Age:         req.Age,
		PhoneNumber: req.PhoneNumber,
		Habits:      req.Habits,
		Alive:       req.Alive,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

func (s *UserServise) UpdateUser(ctx context.Context, req *api.UpdateUserRequest) (*api.UpdateUserResponse, error) {
	var user api.UpdateUserResponse

	err := s.db.QueryRowContext(ctx, `
			UPDATE users
			SET full_name = $2, age = $3, phone_number = $4, habits = $5, alive = $6
			WHERE id = $1
			RETURNING id, full_name, age, phone_number, habits, alive, created_at
    `, req.Id, req.FullName, req.Age, req.PhoneNumber, req.Habits, req.Alive).Scan(
		&user.Id,
		&user.FullName,
		&user.Age,
		&user.PhoneNumber,
		&user.Habits,
		&user.Alive,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &user, nil
}

func (s *UserServise) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.GetUserResponse, error) {
	var user api.GetUserResponse

	err := s.db.QueryRowContext(ctx, `
	SELECT full_name, age, phone_number, habits, alive, created_at FROM users WHERE id = $1
	`, req.Id).Scan(&user.FullName, &user.Age, &user.PhoneNumber, &user.Habits, &user.Alive, &user.CreatedAt)

	user.Id = req.Id

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &user, nil
}

func (s *UserServise) GetAllUsers(ctx context.Context, req *api.GetAllUsersRequest) (*api.GetAllUsersResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 2
	}
	offset := (page - 1) * limit

	sqlStatement := `
		SELECT id, full_name, age, phone_number, habits, alive, created_at 
		FROM users
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, sqlStatement, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to query users: %v", err)
	}
	defer rows.Close()

	var users []*api.GetUserResponse

	for rows.Next() {
		var user api.GetUserResponse

		err := rows.Scan(
			&user.Id,
			&user.FullName,
			&user.Age,
			&user.PhoneNumber,
			&user.Habits,
			&user.Alive,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan user row: %v", err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error after iterating rows: %v", err)
	}

	return &api.GetAllUsersResponse{
		Users: users,
	}, nil
}

func (s *UserServise) DeleteUser(ctx context.Context, req *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	var deletedID int

	err := s.db.QueryRowContext(ctx, `
			DELETE FROM users WHERE id = $1
			RETURNING id
    `, req.Id).Scan(&deletedID)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
