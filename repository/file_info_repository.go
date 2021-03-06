package repository

import (
	"database/sql"
	"fmt"
	"proxy-fileserver/models"
)

type FileInfoRepository struct {
	tableName string
	db        *sql.DB
}

func NewFileInfoRepository(db *sql.DB, tableName string) *FileInfoRepository {
	return &FileInfoRepository{
		db:        db,
		tableName: tableName,
	}
}

func (r *FileInfoRepository) Create(model models.FileInfo) error {
	query := fmt.Sprintf("INSERT INTO %s(filepath, last_download_at) VALUES(?, ?, ?)", r.tableName)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(model.FilePath, model.LastDownloadAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *FileInfoRepository) Update(model models.FileInfo) error {
	query := fmt.Sprintf("UPDATE %s SET last_download_at = ? WHERE filepath = ?", r.tableName)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(model.LastDownloadAt, model.FilePath)
	if err != nil {
		return err
	}
	return nil
}

func (r *FileInfoRepository) Delete(id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id =?", r.tableName)
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
