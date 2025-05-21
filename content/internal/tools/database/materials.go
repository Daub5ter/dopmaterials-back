package database

import (
	"content/internal/models"
	"content/internal/models/errs"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
)

func (db Database) GetMaterials(categoryId, limit *uint32, offset uint32) ([]*models.Material, error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	rows, err := db.conn.Query(ctx, selectMaterials,
		categoryId,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []*models.Material

	for rows.Next() {
		var material models.Material
		err = rows.Scan(
			&material.Id,
			&material.CategoryId,
			&material.Name,
			&material.Description,
			&material.PreviewMeta,
			&material.VideoMeta,
			&material.Deleted,
			&material.CreatedAt,
			&material.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		materials = append(materials, &material)
	}

	return materials, nil
}

func (db Database) GetSearchedMaterials(ids []string) ([]*models.Material, error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	rows, err := db.conn.Query(ctx, selectSearchedMaterials, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	materials := make([]*models.Material, 0, len(ids))
	for rows.Next() {
		var material models.Material
		err = rows.Scan(
			&material.Id,
			&material.CategoryId,
			&material.Name,
			&material.Description,
			&material.PreviewMeta,
			&material.VideoMeta,
			&material.Deleted,
			&material.CreatedAt,
			&material.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		materials = append(materials, &material)
	}

	return materials, nil
}

func (db Database) GetMaterial(id uint64) (*models.Material, error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	var material models.Material

	row := db.conn.QueryRow(ctx, selectMaterial, id)

	err := row.Scan(
		&material.Id,
		&material.CategoryId,
		&material.Name,
		&material.Description,
		&material.PreviewMeta,
		&material.VideoMeta,
		&material.Deleted,
		&material.CreatedAt,
		&material.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &material, nil
}

func (db Database) GetMaterialsIdsByCategory(categoryId uint32) (materialsIds []uint64, err error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	rows, err := db.conn.Query(ctx, selectMaterialsIdsByCategoryId, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id uint64
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		materialsIds = append(materialsIds, id)
	}

	return materialsIds, nil
}

func (db Database) InsertMaterial(material *models.Material) (operationTransaction any, id uint64, err error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadUncommitted})
	if err != nil {
		return nil, 0, err
	}

	err = tx.QueryRow(ctx, insertMaterial,
		material.CategoryId,
		material.Name,
		material.Description,
		material.PreviewMeta,
		material.VideoMeta,
	).Scan(&material.Id)
	if err != nil {
		return nil, 0, err
	}

	return tx, material.Id, nil
}

func (db Database) UpdateMaterial(material *models.Material) (operationTransaction any, err error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadUncommitted})
	if err != nil {
		return nil, err
	}

	_, err = db.conn.Exec(ctx, updateMaterial,
		material.CategoryId,
		material.Name,
		material.Description,
		material.Id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return tx, nil
}

func (db Database) DeleteMaterial(id uint64) (operationTransaction any, err error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadUncommitted})
	if err != nil {
		return nil, err
	}

	_, err = db.conn.Exec(ctx, deleteMaterial, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return tx, nil
}

func (db Database) RestoreMaterial(id uint64) (operationTransaction any, err error) {
	ctx, cancel := context.WithTimeout(db.ctx, db.timeout)
	defer cancel()

	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadUncommitted})
	if err != nil {
		return nil, err
	}

	_, err = db.conn.Exec(ctx, restoreMaterial, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return tx, nil
}
