package database

// createMaterials, createCategories
// - команды создания таблиц базы данных, если они отсутствуют.
const (
	createMaterials = `
		create table if not exists materials (
        id bigserial primary key,
        category_id integer references categories(id),
        name varchar(255) not null,
        description text not null,
        preview_meta varchar(255) not null,
        video_meta varchar(255) not null,
        deleted bool not null default false,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`

	createCategories = `
		create table if not exists categories (
		id serial primary key,
    	category_id integer,
		name varchar(255) not null)`
)

const (
	selectMaterials = `
	select
	id,
	category_id,
	name,
	description,
	preview_meta,
	video_meta,
	deleted,
	created_at,
	updated_at
	from materials
	where ($1::integer is null or category_id = $1)
	order by created_at desc
	limit $2
	offset $3 * 3`

	selectSearchedMaterials = `
		select 
		id,
		category_id,
		name,
		description,
		preview_meta,
		video_meta,
		deleted,
		created_at,
		updated_at
		from materials
		where id = ANY ($1)`

	selectMaterial = `
		select 
		id,
		category_id,
		name,
		description,
		preview_meta,
		video_meta,
		deleted,
		created_at,
		updated_at
		from materials
		where id = $1`

	selectMaterialsIdsByCategoryId = `
		select 
		id,
		from materials
		where category_id = $1`

	selectNullParentCategories = `
		select
		id,
		category_id,
		name
		from categories
		where category_id IS NULL`

	selectSubsidiariesCategories = `
		select
		id,
		category_id,
		name
		from categories
		where category_id = $1`

	selectCategory = `
		select
		id,
		category_id,
		name
		from categories
		where id = $1`
)

const (
	insertMaterial = `
		insert into materials(
		category_id,
	    name,
		description,
	    preview_meta,
	    video_meta) 
		values ($1, $2, $3, $4, $5, $6)
		returning id`

	insertCategory = `
		insert into categories(
		category_id,
		name)
		values($1, $2)
		returning id`
)

const (
	updateMaterial = `
		update materials set 
		category_id = $1,
		name = $2,
		description = $3,
		updated_at = now()
		where id = $4`

	updateCategory = `
		update categories set
		category_id = $1,
		name = $2
		where id = $3`

	updateMaterialCategoryIdNull = `
		update materials set
		category_id = NULL,
		updated_at = now()
		where category_id = $1`
)

const (
	deleteMaterial = `
		update materials set
		deleted = true,
		updated_at = now()
		where id = $1`

	deleteCategory = `
		delete from categories
		where id = $1`
)

const (
	restoreMaterial = `
		update materials set
		deleted = false,
		updated_at = now()
		where id = $1`
)
