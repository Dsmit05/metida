package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Dsmit05/metida/internal/logger"
	"github.com/Dsmit05/metida/internal/models"
	"github.com/Dsmit05/metida/internal/repositories/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

var (
	errUserNotFound    = errors.New("User Not Found")
	errUserIsExist     = errors.New("User already exists")
	errBlogNotFound    = errors.New("Blog Not Found")
	errBlogIsExist     = errors.New("Blog already exists")
	errContentNotFound = errors.New("Content Not Found")
	errContentIsExist  = errors.New("Content already exists")
	errOther           = errors.New("pls, Try again")
)

type DBConnectI interface {
	GetConnectDB() string
}

type metricI interface {
	IncDbError()
}

// PostgresRepository repository for accessing Postgres database
type PostgresRepository struct {
	conn    *pgx.Conn
	queries *postgres.Queries
	metric  metricI
}

func NewPostgresRepository(url DBConnectI, metric metricI) (*PostgresRepository, error) {
	conn, err := pgx.Connect(context.Background(), url.GetConnectDB())
	if err != nil {
		return nil, err
	}
	queries := postgres.New(conn)

	logger.Info("repositories.NewPostgresRepository()", "Init")
	return &PostgresRepository{
		conn:    conn,
		queries: queries,
		metric:  metric,
	}, nil
}

func (o *PostgresRepository) Close() {
	err := o.conn.Close(context.Background())
	if err != nil {
		logger.Error("(o *PostgresRepository) Close() error:", err)
	}
}

func (o *PostgresRepository) CreateUser(name, password, email, role string) error {
	inputData := postgres.CreateUserParams{
		Name:     sql.NullString{name, true},
		Password: password,
		Email:    email,
		Role:     role,
	}

	err := o.queries.CreateUser(context.Background(), inputData)
	val, ok := err.(*pgconn.PgError)

	if ok && pgerrcode.IsIntegrityConstraintViolation(val.Code) {
		return errUserIsExist
	}

	if err != nil {
		logger.DatabaseError("queries.CreateUser", err, inputData)
		return errOther
	}

	return nil
}

func (o *PostgresRepository) ReadUser(email string) (*models.User, error) {
	user, err := o.queries.ReadUser(context.Background(), email)

	if err != nil {
		logger.DatabaseError("queries.ReadUser", err, email)
		if user.Email == "" {
			return nil, errUserNotFound
		}

		return nil, errOther
	}

	if user.IsDeleted.Bool {
		return nil, errUserNotFound
	}

	userModel := &models.User{
		ID:        user.ID,
		Name:      user.Name.String,
		Password:  user.Password,
		Email:     user.Email,
		Role:      user.Role,
		IsDeleted: user.IsDeleted.Bool,
	}

	return userModel, nil
}

func (o *PostgresRepository) UpdateUser(email string, name, password, role string, isDeleted bool) error {
	inputData := postgres.UpdateUserParams{
		Email:     email,
		Name:      sql.NullString{name, true},
		Password:  password,
		Role:      role,
		IsDeleted: sql.NullBool{isDeleted, true},
	}

	err := o.queries.UpdateUser(context.Background(), inputData)
	if err != nil {
		logger.DatabaseError("queries.UpdateUser", err, inputData)
		return errOther
	}

	return nil
}

func (o *PostgresRepository) DeleteUser(email string) error {
	err := o.queries.DeleteUser(context.Background(), email)
	if err != nil {
		logger.DatabaseError("queries.DeleteUser", err, email)
		return errOther
	}

	return nil
}

func (o *PostgresRepository) CreateSession(email string, refreshToken, userAgent, ip string, expiresIn int64) error {
	inputData := postgres.CreateSessionParams{
		UserEmail:    sql.NullString{email, true},
		RefreshToken: sql.NullString{refreshToken, true},
		AccessToken:  sql.NullString{Valid: false},
		UserAgent:    sql.NullString{userAgent, true},
		Ip:           sql.NullString{ip, true},
		ExpiresIn:    expiresIn,
	}

	err := o.queries.CreateSession(context.Background(), inputData)

	val, ok := err.(*pgconn.PgError)
	if ok && pgerrcode.IsIntegrityConstraintViolation(val.Code) {
		return errOther
	}

	if err != nil {
		logger.DatabaseError("queries.CreateSession", err, inputData)
		return errOther
	}

	return nil
}

func (o *PostgresRepository) ReadSession(email string, userAgent, ip string) (*models.Session, error) {
	inputData := postgres.ReadSessionParams{
		UserEmail: sql.NullString{email, true},
		UserAgent: sql.NullString{userAgent, true},
		Ip:        sql.NullString{ip, true},
	}

	session, err := o.queries.ReadSession(context.Background(), inputData)

	if err != nil {
		logger.DatabaseError("queries.ReadSession", err, inputData)
		if session.UserEmail.String == "" {
			return nil, errOther
		}

		return nil, err
	}

	if !session.UserEmail.Valid {
		return nil, errOther
	}

	userModel := &models.Session{
		ID:           session.ID,
		UserEmail:    session.UserEmail.String,
		RefreshToken: session.RefreshToken.String,
		AccessToken:  session.AccessToken.String,
		UserAgent:    session.UserAgent.String,
		IP:           session.Ip.String,
		ExpiresIn:    session.ExpiresIn,
		CreatedAt:    session.CreatedAt,
	}

	return userModel, nil
}

func (o *PostgresRepository) UpdateSession(
	email string, refreshToken string, newRefreshToken string, expiresIn int64) error {
	inputData := postgres.UpdateSessionParams{
		UserEmail:      sql.NullString{email, true},
		RefreshToken:   sql.NullString{refreshToken, true},
		RefreshToken_2: sql.NullString{newRefreshToken, true},
		ExpiresIn:      expiresIn,
	}

	err := o.queries.UpdateSession(context.Background(), inputData)
	if err != nil {
		logger.DatabaseError("queries.UpdateSession", err, inputData)
		return errOther
	}

	return nil
}

func (o *PostgresRepository) UpdateSessionTokenOnly(
	refreshToken string, newRefreshToken string, expiresIn int64) error {

	inputData := postgres.UpdateSessionTokenOnlyParams{
		RefreshToken:   sql.NullString{refreshToken, true},
		RefreshToken_2: sql.NullString{newRefreshToken, true},
		ExpiresIn:      expiresIn,
	}

	err := o.queries.UpdateSessionTokenOnly(context.Background(), inputData)
	if err != nil {
		logger.DatabaseError("queries.UpdateSessionTokenOnly", err, inputData)
		return errOther
	}

	return nil
}

func (o *PostgresRepository) ReadEmailRoleWithRefreshToken(refreshToken string) (*models.UserEmailRole, error) {
	inputData := sql.NullString{
		String: refreshToken,
		Valid:  true,
	}

	emailAndRole, err := o.queries.ReadEmailRoleFromSessions(context.Background(), inputData)

	if err != nil {
		logger.DatabaseError("queries.ReadEmailRoleFromSessions", err, inputData)
		if emailAndRole.Email == "" {
			return nil, errUserNotFound
		}

		return nil, errOther
	}

	userModel := &models.UserEmailRole{
		Email:     emailAndRole.Email,
		Role:      emailAndRole.Role,
		ExpiresIn: emailAndRole.ExpiresIn,
	}

	return userModel, nil
}

func (o *PostgresRepository) DeleteSession(email string, ip, userAgent string) error {
	inputData := postgres.DeleteSessionParams{
		UserAgent: sql.NullString{email, true},
		Ip:        sql.NullString{ip, true},
		UserEmail: sql.NullString{userAgent, true},
	}

	err := o.queries.DeleteSession(context.Background(), inputData)
	if err != nil {
		logger.DatabaseError("queries.DeleteSession", err, inputData)
		return errOther
	}

	return nil
}

func (o *PostgresRepository) CreatContent(email string, name, description string) error {
	inputData := postgres.CreateContentParams{
		UserEmail:   sql.NullString{email, true},
		Name:        sql.NullString{name, true},
		Description: sql.NullString{description, true},
	}

	err := o.queries.CreateContent(context.Background(), inputData)

	val, ok := err.(*pgconn.PgError)
	if ok && pgerrcode.IsIntegrityConstraintViolation(val.Code) {
		return errContentIsExist
	}

	if err != nil {
		logger.DatabaseError("queries.CreateContent", err, inputData)
		return errOther
	}

	return nil
}

func (o *PostgresRepository) ReadContent(email string, id int32) (*models.Content, error) {
	inputData := postgres.ReadContentParams{
		UserEmail: sql.NullString{email, true},
		ID:        id,
	}

	content, err := o.queries.ReadContent(context.Background(), inputData)
	if err != nil {
		logger.DatabaseError("queries.ReadContent", err, inputData)
		return nil, errContentNotFound
	}

	contentModel := &models.Content{
		ID:          content.ID,
		UserEmail:   content.UserEmail.String,
		Name:        content.Name.String,
		Description: content.Description.String,
	}

	return contentModel, nil
}

func (o *PostgresRepository) CreatBlog(name, description string) error {
	inputData := postgres.CreateBlogParams{
		Name:        sql.NullString{name, true},
		Description: sql.NullString{description, true},
	}

	err := o.queries.CreateBlog(context.Background(), inputData)

	val, ok := err.(*pgconn.PgError)
	if ok && pgerrcode.IsIntegrityConstraintViolation(val.Code) {
		return errBlogIsExist
	}

	if err != nil {
		logger.DatabaseError("queries.CreateBlog", err, inputData)
		return errOther
	}

	return nil
}

func (o *PostgresRepository) ReadBlog(id int32) (*models.Blog, error) {
	blog, err := o.queries.ReadBlog(context.Background(), id)
	if err != nil {
		logger.DatabaseError("queries.CreateBlog", err, id)

		// Todo: данная метрика используется пока только тут
		o.metric.IncDbError()
		if blog.Name.String == "" {
			return nil, errBlogNotFound
		}

		return nil, errOther
	}

	if !blog.Description.Valid {
		return nil, errOther
	}

	blogModel := &models.Blog{
		ID:          blog.ID,
		Name:        blog.Name.String,
		Description: blog.Description.String,
	}

	return blogModel, nil
}
