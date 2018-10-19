package auditlog

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"time"

	"github.com/pkg/errors"

	"github.com/fabric8-services/fabric8-common/convert"

	uuid "github.com/satori/go.uuid"
)

// Record an audit log record to track usage of the service
type Record struct {
	ID          uuid.UUID   `sql:"type:uuid default uuid_generate_v4()" gorm:"primary_key"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	IdentityID  uuid.UUID   `sql:"type:uuid"` // TODO: or should we store the sub claim instead as a string: no need to convert it into an UUID each time
	EventTypeID uuid.UUID   `sql:"type:uuid"`
	EventParams EventParams `sql:"type:jsonb"`
}

const (
	recordTableName = "audit_logs"
)

// TableName implements gorm.tabler
func (r Record) TableName() string {
	return recordTableName
}

// Ensure Record implements the Equaler interface
var _ convert.Equaler = Record{}
var _ convert.Equaler = (*Record)(nil)

// Equal returns true if two Record objects are equal; otherwise false is returned.
func (r Record) Equal(o convert.Equaler) bool {
	other, ok := o.(Record)
	if !ok {
		return false
	}
	return r.ID == other.ID
}

// EventParams the parameters of the recorded/audited event
type EventParams map[string]interface{}

// Ensure Fields implements the Equaler interface
var _ convert.Equaler = EventParams{}
var _ convert.Equaler = (*EventParams)(nil)

// Ensure Fields implements the sql.Scanner and driver.Valuer interfaces
var _ sql.Scanner = (*EventParams)(nil)
var _ driver.Valuer = (*EventParams)(nil)

// Equal returns true if two Fields objects are equal; otherwise false is returned.
func (p EventParams) Equal(u convert.Equaler) bool {
	other, ok := u.(EventParams)
	if !ok {
		return false
	}
	return reflect.DeepEqual(p, other)
}

// Value implements the driver.Valuer interface
func (p EventParams) Value() (driver.Value, error) {
	return toBytes(p)
}

// Scan implements the https://golang.org/pkg/database/sql/#Scanner interface
// See also https://stackoverflow.com/a/25374979/835098
// See also https://github.com/jinzhu/gorm/issues/302#issuecomment-80566841
func (p *EventParams) Scan(src interface{}) error {
	return fromBytes(src, p)
}

func toBytes(j interface{}) (driver.Value, error) {
	if j == nil {
		// log.Trace("returning null")
		return nil, nil
	}

	res, error := json.Marshal(j)
	return res, error
}

func fromBytes(src interface{}, target interface{}) error {
	if src == nil {
		target = nil
		return nil
	}
	s, ok := src.([]byte)
	if !ok {
		return errors.Errorf("scanned source was not a string")
	}
	return json.Unmarshal(s, target)
}
