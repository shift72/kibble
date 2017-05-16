package models

// CollectionCollection - all collections
type CollectionCollection []Collection

// Collection - collection of films and tv seasons / episodes
type Collection struct {
	ID          int
	Slug        string
	Title       string
	TitleSlug   string
	Description string
	DisplayName string
	ItemLayout  string
	ItemsPerRow int
	Images      ImageSet
	Seo         Seo
	Items       []GenericItem
	SearchQuery string
	CreatedAt   string
	UpdatedAt   string
}

// FindCollectionByID - find the page by id
func (collections CollectionCollection) FindCollectionByID(collectionID int) (*Collection, error) {
	for _, b := range collections {
		if b.ID == collectionID {
			return &b, nil
		}
	}
	return nil, ErrDataSourceMissing
}

// FindCollectionBySlug - find the collection by the slug
func (collections CollectionCollection) FindCollectionBySlug(slug string) (*Collection, error) {
	for _, p := range collections {
		if p.Slug == slug || p.TitleSlug == slug {
			return &p, nil
		}
	}
	return nil, ErrDataSourceMissing
}

// GetGenericItem - returns a generic item
func (collection Collection) GetGenericItem() GenericItem {
	return GenericItem{
		Title:     collection.Title,
		Seo:       collection.Seo,
		Images:    collection.Images,
		InnerItem: collection,
	}
}
