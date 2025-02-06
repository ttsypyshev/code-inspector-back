package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

// статусы проекта
type Status string

const (
	Draft     Status = "draft"
	Deleted   Status = "deleted"
	Created   Status = "created"
	Completed Status = "completed"
	Rejected  Status = "rejected"
)

func (s Status) IsValid() bool {
	switch s {
	case Draft, Deleted, Created, Completed, Rejected:
		return true
	}
	return false
}

func (s Status) String() string {
	return string(s)
}

// роли для юзеров
type Role string

const (
	Admin   Role = "admin"
	Student Role = "student"
	None Role = "none"
)

func (r Role) IsValid() bool {
	switch r {
	case Admin, Student:
		return true
	}
	return false
}

func (r Role) String() string {
	return string(r)
}

// мапа для списка в языках {key:value}
type JSONB map[string]string

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	var objMap map[string]string
	if err := json.Unmarshal(bytes, &objMap); err == nil {
		*j = objMap
		return nil
	}

	var arr []interface{}
	if err := json.Unmarshal(bytes, &arr); err == nil {
		objMap := make(map[string]string)
		for _, val := range arr {
			if strVal, ok := val.(string); ok {
				parts := strings.SplitN(strVal, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					objMap[key] = value
				}
			} else {
				return errors.New("array elements must be strings in the form 'key:value'")
			}
		}
		*j = objMap
		return nil
	}

	return errors.New("failed to unmarshal JSONB data as either an object or an array")
}
