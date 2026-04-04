package admin

import (
	"context"
	"database/sql"
	"errors"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
	"github.com/martketplace-vkr/order/domain"
)

type repository struct {
	ctxGetter *trmsqlx.CtxGetter
	db        *sqlx.DB
}

func New(db *sqlx.DB, ctxGetter *trmsqlx.CtxGetter) *repository {
	return &repository{
		db:        db,
		ctxGetter: ctxGetter,
	}
}

func (r *repository) InsertUser(ctx context.Context, user *domain.User) error {
	query := `
		insert into employee."user"(
			email,
			password_hash
		) values (
			$1,
			$2
		) returning *
	`

	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		user,
		query,
		user.Email,
		user.PasswordHash,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) SelectUserByEmail(ctx context.Context, email string) (user *domain.User, err error) {
	query := `
		select
			id,
			email,
			password_hash,
			email_verified,
			status,
			created_at,
			updated_at
		from employee."user"
		where email = $1
	`

	user = new(domain.User)

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		user,
		query,
		email,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) HasActiveInviteToken(ctx context.Context, token string) (exists bool, err error) {
	query := `
		select exists(
			select 1
			from employee.invite_token
			where token = $1
				and used_at is null
		)
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&exists,
		query,
		token,
	)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *repository) UseInviteToken(ctx context.Context, token string, userID int64) error {
	query := `
		update employee.invite_token
		set
			used_by = $2,
			used_at = $3
		where token = $1
			and used_at is null
	`

	result, err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
		ctx,
		query,
		token,
		userID,
		time.Now(),
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repository) CreateInviteToken(
	ctx context.Context,
	createdBy int64,
	roleID int64,
	token string,
) error {
	query := `
		insert into employee.invite_token(
			token,
			created_by,
			role_id
		) values (
			$1,
			$2,
			$3
		)
	`

	_, err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
		ctx,
		query,
		token,
		createdBy,
		roleID,
	)
	if err != nil {
		return err
	}

	return nil
}

func IsInviteNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
