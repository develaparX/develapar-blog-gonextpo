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

type CommentRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo CommentRepository
}

func (suite *CommentRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.repo = NewCommentRepository(suite.db)
}

func (suite *CommentRepositoryTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestCreateComment() {
	comment := model.Comment{
		Article: model.Article{Id: 1},
		User:    model.User{Id: 1},
		Content: "Test Comment",
	}

	query := `INSERT INTO comments (article_id, user_id , content, created_at) VALUES($1, $2, $3, $4) RETURNING id, article_id, user_id, content, created_at`

	rows := sqlmock.NewRows([]string{"id", "article_id", "user_id", "content", "created_at"}).
		AddRow(1, comment.Article.Id, comment.User.Id, comment.Content, time.Now())

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(comment.Article.Id, comment.User.Id, comment.Content, sqlmock.AnyArg()).
		WillReturnRows(rows)

	createdComment, err := suite.repo.CreateComment(comment)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), comment.Content, createdComment.Content)
}

func (suite *CommentRepositoryTestSuite) TestCreateComment_Error() {
	comment := model.Comment{
		Article: model.Article{Id: 1},
		User:    model.User{Id: 1},
		Content: "Test Comment",
	}

	query := `INSERT INTO comments (article_id, user_id , content, created_at) VALUES($1, $2, $3, $4) RETURNING id, article_id, user_id, content, created_at`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(comment.Article.Id, comment.User.Id, comment.Content, sqlmock.AnyArg()).
		WillReturnError(errors.New("db error"))

	createdComment, err := suite.repo.CreateComment(comment)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Comment{}, createdComment)
}

func (suite *CommentRepositoryTestSuite) TestGetCommentByArticleId() {
	articleId := 1
	query := `SELECT c.id, c.article_id, c.user_id, c.content, c.created_at, a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, u.created_at, u.updated_at, ca.id, ca.name FROM comments c JOIN articles a ON c.article_id = a.id JOIN users u ON c.user_id = u.id JOIN categories ca ON a.category_id = ca.id WHERE c.article_id = $1 ORDER BY c.created_at DESC`

	rows := sqlmock.NewRows([]string{"id", "article_id", "user_id", "content", "created_at", "article_id_alias", "title", "slug", "content_alias", "views", "article_created_at", "article_updated_at", "user_id_alias", "user_name", "user_email", "user_created_at", "user_updated_at", "category_id", "category_name"}).
		AddRow(1, 1, 1, "Comment 1", time.Now(), 1, "Article Title", "article-slug", "Article Content", 10, time.Now(), time.Now(), 1, "User Name", "user@example.com", time.Now(), time.Now(), 1, "Category Name")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(articleId).WillReturnRows(rows)

	comments, err := suite.repo.GetCommentByArticleId(articleId)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), comments, 1)
	assert.Equal(suite.T(), "Comment 1", comments[0].Content)
}

func (suite *CommentRepositoryTestSuite) TestGetCommentByArticleId_Error() {
	articleId := 1
	query := `SELECT c.id, c.article_id, c.user_id, c.content, c.created_at, a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, u.created_at, u.updated_at, ca.id, ca.name FROM comments c JOIN articles a ON c.article_id = a.id JOIN users u ON c.user_id = u.id JOIN categories ca ON a.category_id = ca.id WHERE c.article_id = $1 ORDER BY c.created_at DESC`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(articleId).WillReturnError(errors.New("db error"))

	comments, err := suite.repo.GetCommentByArticleId(articleId)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), comments)
}

func (suite *CommentRepositoryTestSuite) TestGetCommentByUserId() {
	userId := 1
	query := `SELECT c.id, c.content, c.created_at, u.id, u.name, u.email, a.id, a.title, a.slug FROM comments c JOIN users u ON c.user_id = u.id JOIN articles a ON c.article_id = a.id WHERE c.user_id = $1`

	rows := sqlmock.NewRows([]string{"id", "content", "created_at", "user_id", "user_name", "user_email", "article_id", "article_title", "article_slug"}).
		AddRow(1, "Comment 1", time.Now(), 1, "User Name", "user@example.com", 1, "Article Title", "article-slug")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userId).WillReturnRows(rows)

	comments, err := suite.repo.GetCommentByUserId(userId)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), comments, 1)
	assert.Equal(suite.T(), "Comment 1", comments[0].Content)
}

func (suite *CommentRepositoryTestSuite) TestGetCommentByUserId_Error() {
	userId := 1
	query := `SELECT c.id, c.content, c.created_at, u.id, u.name, u.email, a.id, a.title, a.slug FROM comments c JOIN users u ON c.user_id = u.id JOIN articles a ON c.article_id = a.id WHERE c.user_id = $1`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userId).WillReturnError(errors.New("db error"))

	comments, err := suite.repo.GetCommentByUserId(userId)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), comments)
}

func (suite *CommentRepositoryTestSuite) TestGetCommentById() {
	commentId := 1
	query := `SELECT id, article_id, user_id, content, created_at FROM comments WHERE id = $1`

	rows := sqlmock.NewRows([]string{"id", "article_id", "user_id", "content", "created_at"}).
		AddRow(1, 1, 1, "Test Comment", time.Now())

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(commentId).WillReturnRows(rows)

	comment, err := suite.repo.GetCommentById(commentId)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Comment", comment.Content)
}

func (suite *CommentRepositoryTestSuite) TestGetCommentById_Error() {
	commentId := 1
	query := `SELECT id, article_id, user_id, content, created_at FROM comments WHERE id = $1`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(commentId).WillReturnError(errors.New("db error"))

	comment, err := suite.repo.GetCommentById(commentId)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Comment{}, comment)
}

func (suite *CommentRepositoryTestSuite) TestUpdateComment() {
	commentId := 1
	content := "Updated Comment"
	userId := 1
	query := `UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2 AND user_id=$3`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(content, commentId, userId).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.UpdateComment(commentId, content, userId)
	assert.NoError(suite.T(), err)
}

func (suite *CommentRepositoryTestSuite) TestUpdateComment_Error() {
	commentId := 1
	content := "Updated Comment"
	userId := 1
	query := `UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2 AND user_id=$3`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(content, commentId, userId).WillReturnError(errors.New("db error"))

	err := suite.repo.UpdateComment(commentId, content, userId)
	assert.Error(suite.T(), err)
}

func (suite *CommentRepositoryTestSuite) TestDeleteComment() {
	commentId := 1
	query := `DELETE FROM comments WHERE id = $1`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(commentId).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.DeleteComment(commentId)
	assert.NoError(suite.T(), err)
}

func (suite *CommentRepositoryTestSuite) TestDeleteComment_Error() {
	commentId := 1
	query := `DELETE FROM comments WHERE id = $1`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(commentId).WillReturnError(errors.New("db error"))

	err := suite.repo.DeleteComment(commentId)
	assert.Error(suite.T(), err)
}

func TestCommentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CommentRepositoryTestSuite))
}
