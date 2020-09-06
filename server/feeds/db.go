package feeds

import (
	"context"
	"database/sql"
)

type feedsDB struct {
	db *sql.DB
}

func newDB(sqldb *sql.DB) *feedsDB {
	return &feedsDB{
		db: sqldb,
	}
}

// Create adds a new articles to the database.
func (db *feedsDB) Create(ctx context.Context, article *Feed) error {
	return nil
	//createTime, err := ptypes.Timestamp(article.CreateTime)
	//if err != nil {
	//return fmt.Errorf("failed to convret create_time: %w", err)
	//}

	//_, createErr := db.db.ExecContext(ctx, `
	//INSERT INTO articles (
	//id,
	//user_id,
	//url,
	//title,
	//create_time,
	//content,
	//content_sha256,
	//is_read,
	//source_id
	//) VALUES (
	//$1, $2, $3, $4, $5, $6, $7, $8, $9
	//)
	//`,
	//article.Id,
	//article.UserId,
	//article.Url,
	//article.Title,
	//createTime.UnixNano(),
	//article.Content,
	//article.ContentSha256,
	//article.IsRead,
	//article.SourceId,
	//)
	//return createErr
}

// Get returns article by id.
func (db *feedsDB) Get(ctx context.Context, id string) (*Feed, error) {
	return nil, nil
	//row := db.db.QueryRowContext(ctx, `
	//SELECT
	//id,
	//user_id,
	//url,
	//title,
	//create_time,
	//content,
	//content_sha256,
	//is_read,
	//source_id
	//FROM
	//articles
	//WHERE
	//id = $1
	//`, id)

	//return db.scanRow(row)
}

//func (db *feedsDB) scanRow(row *sql.Row) (*Feed, error) {
//err := row.Scan(
//&article.Article.Id,
//&article.Article.UserId,
//&article.Article.Url,
//&article.Article.Title,
//&article.CreateTimestamp,
//&article.Article.Content,
//&article.Article.ContentSha256,
//&article.Article.IsRead,
//&article.Article.SourceId,
//)

//if err != nil {
//return nil, err
//}

//var convertTimeErr error
//article.Article.CreateTime, convertTimeErr = ptypes.TimestampProto(time.Unix(0, article.CreateTimestamp))
//if convertTimeErr != nil {
//return nil, fmt.Errorf("failed to convert create time: %w", convertTimeErr)
//}

//return article.Article, nil
//}

// Delete deletes an article by id.
func (db *feedsDB) Delete(ctx context.Context, id string) error {
	return nil
	//_, err := db.db.ExecContext(ctx, `
	//DELETE FROM feeds
	//WHERE id = $1
	//`, id)
	//return err
}

func (db *feedsDB) Update(ctx context.Context, article *Article) error {
	return nil
	//_, err := db.db.ExecContext(ctx, `
	//UPDATE articles
	//SET
	//is_read = $1,
	//content = $2,
	//content_sha256 = $3
	//WHERE
	//id = $4
	//`, article.IsRead, article.Content, article.ContentSha256, article.Id)
	//return err
}

// List returns all articles.
func (db *feedsDB) List(ctx context.Context, request *ListArticlesRequest) ([]*Article, error) {
	return nil, nil
	//q := &strings.Builder{}
	//q.WriteString(`
	//SELECT
	//id,
	//user_id,
	//url,
	//title,
	//create_time,
	//content_sha256,
	//is_read,
	//source_id
	//FROM articles
	//WHERE
	//user_id = $1
	//AND id >= $2
	//`)

	//from, err := request.FromID()
	//if err != nil {
	//return nil, err
	//}

	//args := []interface{}{
	//request.UserId, from,
	//}

	//if request.IsReadEq != nil {
	//args = append(args, request.GetIsReadEq().GetValue())
	//q.WriteString(fmt.Sprintf(`
	//AND is_read = $%d
	//`, len(args)))
	//}

	//if request.GetSourceIdEq() != nil {
	//args = append(args, request.GetSourceIdEq().GetValue())
	//q.WriteString(fmt.Sprintf(`
	//AND source_id = $%d
	//`, len(args)))
	//}

	//if request.GetTitleContains() != nil {
	//args = append(args, fmt.Sprintf("%%%s%%", request.GetTitleContains().GetValue()))
	//q.WriteString(fmt.Sprintf(`
	//AND title LIKE $%d
	//`, len(args)))
	//}

	//args = append(args, request.PageSize)
	//q.WriteString(fmt.Sprintf(`
	//ORDER BY id ASC
	//LIMIT $%d
	//`, len(args)))

	//rows, err := db.db.QueryContext(ctx, q.String(), args...)
	//if err != nil {
	//return nil, err
	//}
	//defer rows.Close()
	//articles := []*Article{}
	//for rows.Next() {
	//article := &dbArticle{
	//Article: &Article{},
	//}

	//err := rows.Scan(
	//&article.Article.Id,
	//&article.Article.UserId,
	//&article.Article.Url,
	//&article.Article.Title,
	//&article.CreateTimestamp,
	//&article.Article.ContentSha256,
	//&article.Article.IsRead,
	//&article.Article.SourceId,
	//)

	//if err != nil {
	//return nil, err
	//}

	//var convertTimeErr error
	//article.Article.CreateTime, convertTimeErr = ptypes.TimestampProto(time.Unix(0, article.CreateTimestamp))
	//if convertTimeErr != nil {
	//return nil, fmt.Errorf("failed to convert create time: %w", convertTimeErr)
	//}

	//articles = append(articles, article.Article)
	//}
	//if err := rows.Close(); err != nil {
	//return nil, err
	//}
	//if err := rows.Err(); err != nil {
	//return nil, err
	//}
	//return articles, nil
}
