package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TagRepository interface {
	Create(ctx context.Context, tag *model.Tag) error
	CreateBatch(ctx context.Context, tags []*model.Tag) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Tag, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Tag, error)
	GetByName(ctx context.Context, name string) (*model.Tag, error)
	GetByNames(ctx context.Context, names []string) ([]*model.Tag, error)
	List(ctx context.Context, filter TagFilter) ([]*model.Tag, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteBatch(ctx context.Context, ids []uuid.UUID) error
	GetOrCreate(ctx context.Context, name string) (*model.Tag, error)
	GetOrCreateBatch(ctx context.Context, names []string) ([]*model.Tag, error)
}

type TagFilter struct {
	IDs     []uuid.UUID
	Names   []string
	Search  string
	Limit   int
	Offset  int
	OrderBy string
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) getDB(ctx context.Context) *gorm.DB {
	if txWrapper, ok := ctx.Value(ctxKey{}).(*Transaction); ok && txWrapper.Tx != nil {
		return txWrapper.Tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *tagRepository) Create(ctx context.Context, tag *model.Tag) error {
	return r.getDB(ctx).Create(tag).Error
}

func (r *tagRepository) CreateBatch(ctx context.Context, tags []*model.Tag) error {
	if len(tags) == 0 {
		return nil
	}
	return r.getDB(ctx).Create(tags).Error
}

func (r *tagRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Tag, error) {
	var tag model.Tag
	err := r.getDB(ctx).First(&tag, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Tag, error) {
	if len(ids) == 0 {
		return []*model.Tag{}, nil
	}

	var tags []*model.Tag
	err := r.getDB(ctx).Where("id IN ?", ids).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepository) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.getDB(ctx).Where("name = ?", name).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetByNames(ctx context.Context, names []string) ([]*model.Tag, error) {
	if len(names) == 0 {
		return []*model.Tag{}, nil
	}

	var tags []*model.Tag
	err := r.getDB(ctx).Where("name IN ?", names).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepository) List(ctx context.Context, filter TagFilter) ([]*model.Tag, int64, error) {
	query := r.getDB(ctx).Model(&Tag{})

	if len(filter.IDs) > 0 {
		query = query.Where("id IN ?", filter.IDs)
	}
	if len(filter.Names) > 0 {
		query = query.Where("name IN ?", filter.Names)
	}
	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.OrderBy != "" {
		query = query.Order(filter.OrderBy)
	} else {
		query = query.Order("name ASC")
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var tags []*model.Tag
	err := query.Find(&tags).Error
	if err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.getDB(ctx).Delete(&Tag{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *tagRepository) DeleteBatch(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.getDB(ctx).Where("id IN ?", ids).Delete(&Tag{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *tagRepository) GetOrCreate(ctx context.Context, name string) (*model.Tag, error) {
	tag, err := r.GetByName(ctx, name)
	if err == nil {
		return tag, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	tag = &Tag{
		ID:   uuid.New(),
		Name: name,
	}
	if err := r.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (r *tagRepository) GetOrCreateBatch(ctx context.Context, names []string) ([]*model.Tag, error) {
	if len(names) == 0 {
		return []*model.Tag{}, nil
	}

	uniqueNames := compact(names)

	tagsToCreate := make([]*model.Tag, len(uniqueNames))
	for i, name := range uniqueNames {
		tagsToCreate[i] = &model.Tag{
			ID:   uuid.New(),
			Name: name,
		}
	}

	err := r.getDB(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).
		Create(&tagsToCreate).Error

	if err != nil {
		return nil, err
	}

	var tags []*model.Tag
	err = r.getDB(ctx).Where("name IN ?", uniqueNames).Find(&tags).Error
	return tags, err
}

func compact(names []string) []string {
	if len(names) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(names))
	result := make([]string, 0, len(names))

	for _, name := range names {
		if _, exists := seen[name]; !exists {
			seen[name] = struct{}{}
			result = append(result, name)
		}
	}

	return result
}
