package domain

type SystemConfig struct {
    Base
    ConfigKey string `json:"config_key,omitempty"`
    ConfigValue string `json:"config_value,omitempty"`
    ValueType string `json:"value_type,omitempty"`
    Description string `json:"description,omitempty"`
    UpdatedBy *int64 `json:"updated_by,omitempty"`
}

func (m SystemConfig) EntityID() int64 { return m.ID }
func (m SystemConfig) IsDeleted() bool { return m.DeletedAt != nil }
func (m SystemConfig) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m SystemConfig) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }