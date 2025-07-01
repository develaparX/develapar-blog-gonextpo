package repository_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"develapar-server/model"
	"develapar-server/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BookmarkRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo repository.BookmarkRepository
}

func (suite *BookmarkRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.repo = repository.NewBookmarkRepository(suite.db)
}

func (suite *BookmarkRepositoryTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *BookmarkRepositoryTestSuite) TestCreateBookmark() {
	bookmark := model.Bookmark{
		Article: model.Article{Id: 1},
		User:    model.User{Id: 1},
	}

	query := `INSERT INTO bookmarks (article_id, user_id, created_at) VALUES ($1, $2, $3) RETURNING id, article_id, user_id,created_at`

	rows := sqlmock.NewRows([]string{"id", "article_id", "user_id", "created_at"}).
		AddRow(1, bookmark.Article.Id, bookmark.User.Id, time.Now())

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(bookmark.Article.Id, bookmark.User.Id, sqlmock.AnyArg()).
		WillReturnRows(rows)

	createdBookmark, err := suite.repo.CreateBookmark(bookmark)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), bookmark.Article.Id, createdBookmark.Article.Id)
}

func (suite *BookmarkRepositoryTestSuite) TestCreateBookmark_Error() {
	bookmark := model.Bookmark{
		Article: model.Article{Id: 1},
		User:    model.User{Id: 1},
	}

	query := `INSERT INTO bookmarks (article_id, user_id, created_at) VALUES ($1, $2, $3) RETURNING id, article_id, user_id,created_at`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(bookmark.Article.Id, bookmark.User.Id, sqlmock.AnyArg()).
		WillReturnError(errors.New("db error"))

	createdBookmark, err := suite.repo.CreateBookmark(bookmark)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Bookmark{}, createdBookmark)
}

func (suite *BookmarkRepositoryTestSuite) TestGetByUserId() {
	userId := "1"
	query := `SELECT b.id, b.article_id, b.user_id, b.created_at, a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, u.created_at, u.updated_at, c.id, c.name FROM bookmarks b JOIN articles a ON b.article_id = a.id JOIN users u ON b.user_id = u.id JOIN categories c ON a.category_id = c.id WHERE b.user_id = $1 ORDER BY b.created_at DESC;`

	rows := sqlmock.NewRows([]string{"id", "article_id", "user_id", "created_at", "article_id_alias", "title", "slug", "content", "views", "article_created_at", "article_updated_at", "user_id_alias", "user_name", "user_email", "user_created_at", "user_updated_at", "category_id", "category_name"}).
		AddRow(1, 1, 1, time.Now(), 1, "Article Title", "article-slug", "Article Content", 10, time.Now(), time.Now(), 1, "User Name", "user@example.com", time.Now(), time.Now(), 1, "Category Name")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userId).WillReturnRows(rows)

	bookmarks, err := suite.repo.GetByUserId(userId)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), bookmarks, 1)
	assert.Equal(suite.T(), "Article Title", bookmarks[0].Article.Title)
}

func (suite *BookmarkRepositoryTestSuite) TestGetByUserId_Error() {
	userId := "1"
	query := `SELECT b.id, b.article_id, b.user_id, b.created_at, a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, u.created_at, u.updated_at, c.id, c.name FROM bookmarks b JOIN articles a ON b.article_id = a.id JOIN users u ON b.user_id = u.id JOIN categories c ON a.category_id = c.id WHERE b.user_id = $1 ORDER BY b.created_at DESC;`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userId).WillReturnError(errors.New("db error"))

	bookmarks, err := suite.repo.GetByUserId(userId)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), bookmarks)
}

func (suite *BookmarkRepositoryTestSuite) TestDeleteBookmark() {
	userId := 1
	articleId := 1
	query := `DELETE FROM bookmarks WHERE user_id = $1 AND article_id = $2`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(userId, articleId).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.DeleteBookmark(userId, articleId)
	assert.NoError(suite.T(), err)
}

func (suite *BookmarkRepositoryTestSuite) TestDeleteBookmark_Error() {
	userId := 1
	articleId := 1
	query := `DELETE FROM bookmarks WHERE user_id = $1 AND article_id = $2`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(userId, articleId).WillReturnError(errors.New("db error"))

	err := suite.repo.DeleteBookmark(userId, articleId)
	assert.Error(suite.T(), err)
}

func (suite *BookmarkRepositoryTestSuite) TestIsBookmarked() {
	userId := 1
	articleId := 1
	query := `SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = $1 AND article_id = $2)`

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userId, articleId).WillReturnRows(rows)

	isBookmarked, err := suite.repo.IsBookmarked(userId, articleId)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isBookmarked)
}

func (suite *BookmarkRepositoryTestSuite) TestIsBookmarked_Error() {
	userId := 1
	articleId := 1
	query := `SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = $1 AND article_id = $2)`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userId, articleId).WillReturnError(errors.New("db error"))

	isBookmarked, err := suite.repo.IsBookmarked(userId, articleId)
	assert.Error(suite.T(), err)
	assert.False(suite.T(), isBookmarked)
}

func TestBookmarkRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BookmarkRepositoryTestSuite))
}

