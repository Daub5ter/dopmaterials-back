package database

import (
	"context"
	"errors"

	"content/internal/models"
	"content/internal/models/errs"

	"github.com/jackc/pgx/v4"
)

func (db Database) GetNullParentCategories() ([]*models.Category, error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	rows, err := db.conn.Query(ctx, selectNullParentCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category

	for rows.Next() {
		var category models.Category
		err = rows.Scan(
			&category.Id,
			&category.CategoryId,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

func (db Database) GetSubsidiariesCategories(categoryId uint32) ([]*models.Category, error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	rows, err := db.conn.Query(ctx, selectSubsidiariesCategories, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category

	for rows.Next() {
		var category models.Category
		err = rows.Scan(
			&category.Id,
			&category.CategoryId,
			&category.Name,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errs.ErrNotFound
			}
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

func (db Database) GetCategory(id uint32) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	var category models.Category

	row := db.conn.QueryRow(ctx, selectCategory, id)

	err := row.Scan(
		&category.Id,
		&category.CategoryId,
		&category.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &category, nil
}

func (db Database) InsertCategory(category *models.Category) (id uint32, err error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	err = db.conn.QueryRow(ctx, insertCategory,
		category.CategoryId,
		category.Name,
	).Scan(&category.Id)

	return category.Id, nil
}

func (db Database) UpdateCategory(category *models.Category) error {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	_, err := db.conn.Exec(ctx, updateCategory,
		category.CategoryId,
		category.Name,
		category.Id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrNotFound
		}
		return err
	}

	return nil
}

func (db Database) DeleteCategory(id uint32) (operationTransaction any, err error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout*2)
	defer cancel()

	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadUncommitted})
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, updateMaterialCategoryIdNull, id)
	if err != nil {
		_ = tx.Rollback(ctx)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	_, err = tx.Exec(ctx, deleteCategory, id)
	if err != nil {
		_ = tx.Rollback(ctx)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return tx, nil
}
