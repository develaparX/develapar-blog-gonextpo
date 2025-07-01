package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"develapar-server/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LikeRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo LikeRepository
}

func (suite *LikeRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.repo = NewLikeRepository(suite.db)
}

func (suite *LikeRepositoryTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *LikeRepositoryTestSuite) TestCreateLike() {
	like := model.Likes{
		Article: model.Article{Id: 1},
		User:    model.User{Id: 1},
	}

	query := `INSERT INTO likes(article_id, user_id, created_at) VALUES($1,$2,$3) RETURNING id, article_id, user_id, created_at`

	rows := sqlmock.NewRows([]string{"id", "article_id", "user_id", "created_at"}).
		AddRow(1, like.Article.Id, like.User.Id, time.Now())

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(like.Article.Id, like.User.Id, sqlmock.AnyArg()).WillReturnRows(rows)

	createdLike, err := suite.repo.CreateLike(like)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), like.Article.Id, createdLike.Article.Id)
}

func (suite *LikeRepositoryTestSuite) TestCreateLike_Error() {
	like := model.Likes{
		Article: model.Article{Id: 1},
		User:    model.User{Id: 1},
	}

	query := `INSERT INTO likes(article_id, user_id, created_at) VALUES($1,$2,$3) RETURNING id, article_id, user_id, created_at`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(like.Article.Id, like.User.Id, sqlmock.AnyArg()).
		WillReturnError(errors.New("db error"))

	createdLike, err := suite.repo.CreateLike(like)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Likes{}, createdLike)
}

func (suite *LikeRepositoryTestSuite) TestGetLikeByArticleId() {
	articleId := 1
	query := `SELECT l.id, l.article_id, l.user_id, l.created_at, u.id, u.name, u.email FROM likes l JOIN users u ON l.user_id = u.id WHERE l.article_id = $1`

	rows := sqlmock.NewRows([]string{"id", "article_id", "user_id", "created_at", "user_id_alias", "user_name", "user_email"}).
		AddRow(1, 1, 1, time.Now(), 1, "User Name", "user@example.com")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(articleId).WillReturnRows(rows)

	likes, err := suite.repo.GetLikeByArticleId(articleId)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), likes, 1)
	assert.Equal(suite.T(), "User Name", likes[0].User.Name)
}

func (suite *LikeRepositoryTestSuite) TestGetLikeByArticleId_Error() {
	articleId := 1
	query := `SELECT l.id, l.article_id, l.user_id, l.created_at, u.id, u.name, u.email FROM likes l JOIN users u ON l.user_id = u.id WHERE l.article_id = $1`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(articleId).WillReturnError(errors.New("db error"))

	likes, err := suite.repo.GetLikeByArticleId(articleId)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), likes)
}

func (suite *LikeRepositoryTestSuite) TestGetLikeByUserId() {
	userId := 1
	query := `SELECT l.id, l.article_id, l.user_id, l.created_at, a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at FROM likes l JOIN articles a ON l.article_id = a.id WHERE l.user_id = $1`

	rows := sqlmock.NewRows([]string{"id", "article_id", "user_id", "created_at", "article_id_alias", "title", "slug", "content", "views", "article_created_at", "article_updated_at"}).
		AddRow(1, 1, 1, time.Now(), 1, "Article Title", "article-slug", "Article Content", 10, time.Now(), time.Now())

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userId).WillReturnRows(rows)

	likes, err := suite.repo.GetLikeByUserId(userId)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), likes, 1)
	assert.Equal(suite.T(), "Article Title", likes[0].Article.Title)
}

func (suite *LikeRepositoryTestSuite) TestGetLikeByUserId_Error() {
	userId := 1
	query := `SELECT l.id, l.article_id, l.user_id, l.created_at, a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at FROM likes l JOIN articles a ON l.article_id = a.id WHERE l.user_id = $1`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userId).WillReturnError(errors.New("db error"))

	likes, err := suite.repo.GetLikeByUserId(userId)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), likes)
}

func (suite *LikeRepositoryTestSuite) TestDeleteLike() {
	userId := 1
	articleId := 1
	query := `DELETE FROM likes WHERE user_id=$1 AND article_id = $2`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(userId, articleId).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.DeleteLike(userId, articleId)
	assert.NoError(suite.T(), err)
}

func (suite *LikeRepositoryTestSuite) TestDeleteLike_Error() {
	userId := 1
	articleId := 1
	query := `DELETE FROM likes WHERE user_id=$1 AND article_id = $2`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(userId, articleId).WillReturnError(errors.New("db error"))

	err := suite.repo.DeleteLike(userId, articleId)
	assert.Error(suite.T(), err)
}

func (suite *LikeRepositoryTestSuite) TestIsLiked() {
	userId := 1
	articleId := 1
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND article_id = $2)`

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userId, articleId).WillReturnRows(rows)

	isLiked, err := suite.repo.IsLiked(userId, articleId)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isLiked)
}

func (suite *LikeRepositoryTestSuite) TestIsLiked_Error() {
	userId := 1
	articleId := 1
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND article_id = $2)`

	_ = sqlmock.NewRows([]string{"exists"}).AddRow(false)

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userId, articleId).WillReturnError(errors.New("db error"))

	isLiked, err := suite.repo.IsLiked(userId, articleId)
	assert.Error(suite.T(), err)
	assert.False(suite.T(), isLiked)
}

func TestLikeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LikeRepositoryTestSuite))
}
