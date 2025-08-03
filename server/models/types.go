package models

import (
	"database/sql/driver"
	"fmt"
)

// PackageType represents the type of package
type PackageType string

const (
	PackageTypeLibrary  PackageType = "library"
	PackageTypeShowcase PackageType = "showcase"
)

func (pt *PackageType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*pt = PackageType(s)
	case []byte:
		*pt = PackageType(s)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (pt PackageType) Value() (driver.Value, error) {
	return string(pt), nil
}

// PackageStatus represents the status of a package
type PackageStatus string

const (
	PackageStatusInWork    PackageStatus = "in_work"
	PackageStatusReady     PackageStatus = "ready"
	PackageStatusArchived  PackageStatus = "archived"
	PackageStatusAbandoned PackageStatus = "abandoned"
)

func (ps *PackageStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*ps = PackageStatus(s)
	case []byte:
		*ps = PackageStatus(s)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (ps PackageStatus) Value() (driver.Value, error) {
	return string(ps), nil
}

// TagCategory represents the category of a tag
type TagCategory string

const (
	TagCategoryDomain   TagCategory = "domain"
	TagCategoryPlatform TagCategory = "platform"
	TagCategoryFeature  TagCategory = "feature"
)

func (tc *TagCategory) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*tc = TagCategory(s)
	case []byte:
		*tc = TagCategory(s)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (tc TagCategory) Value() (driver.Value, error) {
	return string(tc), nil
}