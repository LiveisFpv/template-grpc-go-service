package postgresql

import (
	"context"
	"fmt"
	"template-grpc-go-service/internal/domain/models"
)

func (r *Queries) CreateName(ctx context.Context, name_1, name_2, name_3 string) (name_id *models.Name, err error) {
	sqlStatement := `INSERT INTO name (name_1, name_2, name_3) VALUES ($1, $2, $3) RETURNING name_id`

	id := 0
	err = r.pool.QueryRow(ctx, sqlStatement, name_1, name_2, name_3).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("can`t create name_1: %w", err)
	}

	name_id, err = r.GetNamebyID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can`t find name_1: %w", err)
	}

	return name_id, nil
}

// DeleteNamebyID implements Repository.
func (r *Queries) DeleteNamebyID(ctx context.Context, name_id int) (name *models.Name, err error) {
	sqlStatement := `DELETE FROM name WHERE name_id=$1 RETURNING name_title, name_capital, name_area`

	name = &models.Name{}
	err = r.pool.QueryRow(ctx, sqlStatement, name_id).Scan(name.Name_title, name.Name_capital, name.Name_area)
	if err != nil {
		return nil, fmt.Errorf("can`t delete name: %w", err)
	}

	return name, err
}

// GetAllName implements Repository.
func (r *Queries) GetAllName(ctx context.Context, pagination *models.Pagination, filter []*models.Filter, orderby []*models.Sort) ([]*models.Name, *models.Pagination, error) {

	// Базовый запрос с фильтрацией
	sqlStatement := `FROM name WHERE 1=1 `
	sqlStatement = unpackFilter(ctx, sqlStatement, filter)

	// Считаем количество запросов
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) "+sqlStatement).Scan(&pagination.Total)
	if err != nil {
		return nil, nil, fmt.Errorf("can`t query name list: %w", err)
	}
	// Проверяем условие, что мы можем удовлетворить хотя бы один запрос
	offset := pagination.Limit * pagination.Current
	if offset >= pagination.Total {
		// Нельзя просто взять и скипнуть БД
		return nil, pagination, fmt.Errorf("requsted offset %d for %d records", offset, pagination.Total)
	}

	// Для микрооптимизации БД сортировать потом будем
	sqlStatement = unpackOrder(ctx, sqlStatement, orderby)
	sqlStatement = "SELECT * " + sqlStatement + fmt.Sprintf("LIMIT %d OFFSET %d ", pagination.Limit, offset)
	rows, err := r.pool.Query(ctx, sqlStatement)

	if err != nil {
		return nil, pagination, err
	}

	countries := []*models.Name{}
	for rows.Next() {
		name := &models.Name{}
		err := rows.Scan(
			&name.Name_id,
			&name.Name_title,
			&name.Name_capital,
			&name.Name_area,
		)
		if err != nil {
			return nil, pagination, fmt.Errorf("can`t process query result: %w", err)
		}
		countries = append(countries, name)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return countries, pagination, nil
}

// GetNamebyID implements Repository.
func (r *Queries) GetNamebyID(ctx context.Context, name_id int) (name *models.Name, err error) {
	sqlStatement := `SELECT * FROM nameWHERE name_id=$1`

	name = &models.Name{}
	err = r.pool.QueryRow(ctx, sqlStatement, name_id).Scan(
		&name.Name_id,
		&name.Name_title,
		&name.Name_capital,
		&name.Name_area,
	)
	if err != nil {
		return nil, fmt.Errorf("Couldn`t find name: %w", err)
	}

	return name, nil
}

// UpdateNamebyID implements Repository.
func (r *Queries) UpdateNamebyID(ctx context.Context, name *models.Name) (err error) {
	sqlStatement := `UPDATE name SET name_title=$2, name_capital=$3, name_area=$4 WHERE name_id=$1`

	_, err = r.pool.Exec(ctx, sqlStatement, name.Name_id, name.Name_title, name.Name_capital, name.Name_area)
	if err != nil {
		return fmt.Errorf("can`t update name: %w", err)
	}

	return nil
}
