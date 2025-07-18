package repository_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"develapar-server/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ArticleTagRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo repository.ArticleTagRepository
}

func (suite *ArticleTagRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.repo = repository.NewArticleTagRepository(suite.db)
}

func (suite *ArticleTagRepositoryTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *ArticleTagRepositoryTestSuite) TestAssignTags() {
	articleId := 1
	tagIds := []int{1, 2}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM article_tags WHERE article_id = $1`)).
		WithArgs(articleId).
		WillReturnResult(sqlmock.NewResult(0, 0))

	suite.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO article_tags (article_id, tag_id) VALUES ($1,$2)`)).
		WithArgs(articleId, tagIds[0]).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO article_tags (article_id, tag_id) VALUES ($1,$2)`)).
		WithArgs(articleId, tagIds[1]).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.mock.ExpectCommit()

	err := suite.repo.AssignTags(articleId, tagIds)
	assert.NoError(suite.T(), err)
}

func (suite *ArticleTagRepositoryTestSuite) TestAssignTags_Rollback() {
	articleId := 1
	tagIds := []int{1, 2}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM article_tags WHERE article_id = $1`)).
		WithArgs(articleId).
		WillReturnResult(sqlmock.NewResult(0, 0))

	suite.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO article_tags (article_id, tag_id) VALUES ($1,$2)`)).
		WithArgs(articleId, tagIds[0]).
		WillReturnError(errors.New("insert error"))

	suite.mock.ExpectRollback()

	err := suite.repo.AssignTags(articleId, tagIds)
	assert.Error(suite.T(), err)
}

func (suite *ArticleTagRepositoryTestSuite) TestGetTagsByArticleId() {
	articleId := 1
	query := `SELECT t.id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.article_id = $1`

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Tag1").
		AddRow(2, "Tag2")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(articleId).WillReturnRows(rows)

	tags, err := suite.repo.GetTagsByArticleId(articleId)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), tags, 2)
	assert.Equal(suite.T(), "Tag1", tags[0].Name)
}

func (suite *ArticleTagRepositoryTestSuite) TestGetTagsByArticleId_Error() {
	articleId := 1
	query := `SELECT t.id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.article_id = $1`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(articleId).WillReturnError(errors.New("db error"))

	tags, err := suite.repo.GetTagsByArticleId(articleId)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), tags)
}

func (suite *ArticleTagRepositoryTestSuite) TestGetArticleByTagId() {
	tagId := 1
	query := `SELECT a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, c.id, c.name FROM articles a JOIN users u ON a.user_id = u.id JOIN categories c ON a.category_id = c.id JOIN article_tags at ON at.article_id = a.id WHERE at.tag_id = $1 ORDER BY a.created_at DESC`

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "views", "created_at", "updated_at", "user_id", "user_name", "user_email", "category_id", "category_name"}).
		AddRow(1, "Article1", "slug1", "Content1", 10, time.Now(), time.Now(), 1, "User1", "user1@example.com", 1, "Category1")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(tagId).WillReturnRows(rows)

	articles, err := suite.repo.GetArticleByTagId(tagId)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), articles, 1)
	assert.Equal(suite.T(), "Article1", articles[0].Title)
}

func (suite *ArticleTagRepositoryTestSuite) TestGetArticleByTagId_Error() {
	tagId := 1
	query := `SELECT a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, c.id, c.name FROM articles a JOIN users u ON a.user_id = u.id JOIN categories c ON a.category_id = c.id JOIN article_tags at ON at.article_id = a.id WHERE at.tag_id = $1 ORDER BY a.created_at DESC`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(tagId).WillReturnError(errors.New("db error"))

	articles, err := suite.repo.GetArticleByTagId(tagId)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), articles)
}

func (suite *ArticleTagRepositoryTestSuite) TestRemoveTagFromArticle() {
	articleId := 1
	tagId := 1
	query := `DELETE FROM article_tags WHERE article_id= $1 AND tag_id = $2`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(articleId, tagId).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.RemoveTagFromArticle(articleId, tagId)
	assert.NoError(suite.T(), err)
}

func (suite *ArticleTagRepositoryTestSuite) TestRemoveTagFromArticle_Error() {
	articleId := 1
	tagId := 1
	query := `DELETE FROM article_tags WHERE article_id= $1 AND tag_id = $2`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(articleId, tagId).WillReturnError(errors.New("db error"))

	err := suite.repo.RemoveTagFromArticle(articleId, tagId)
	assert.Error(suite.T(), err)
}

func TestArticleTagRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleTagRepositoryTestSuite))
}
