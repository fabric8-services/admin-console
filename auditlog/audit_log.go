package auditlog

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"time"

	"github.com/fabric8-services/fabric8-common/closeable"
	"github.com/fabric8-services/fabric8-common/convert"
	"github.com/fabric8-services/fabric8-common/errors"
	"github.com/fabric8-services/fabric8-common/log"

	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	errs "github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// AuditLog an audit log record to track usage of the service
type AuditLog struct {
	ID          uuid.UUID   `sql:"type:uuid default uuid_generate_v4()" gorm:"primary_key;column:audit_log_id"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	IdentityID  uuid.UUID   `sql:"type:uuid"`
	Username    string      `sql:"type:string"`
	EventTypeID uuid.UUID   `sql:"type:uuid"`
	EventParams EventParams `sql:"type:jsonb"`
}

const (
	recordTableName = "audit_log"
)

// TableName implements gorm.tabler
func (r AuditLog) TableName() string {
	return recordTableName
}

// Ensure AuditLog implements the Equaler interface
var _ convert.Equaler = AuditLog{}
var _ convert.Equaler = (*AuditLog)(nil)

// Equal returns true if two AuditLog objects are equal; otherwise false is returned.
func (r AuditLog) Equal(o convert.Equaler) bool {
	other, ok := o.(AuditLog)
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
		return errs.Errorf("scanned source was not a string")
	}
	return json.Unmarshal(s, target)
}

// Repository provides functions to create and view audit logs
type Repository interface {
	Create(ctx context.Context, auditLog *AuditLog) error
	LoadByID(ctx context.Context, id uuid.UUID) (AuditLog, error)
	ListByIdentityID(ctx context.Context, identityID uuid.UUID, start int, limit int) ([]AuditLog, int, error)
	ListByUsername(ctx context.Context, username string, start int, limit int) ([]AuditLog, int, error)
}

// NewRepository creates a GormRecordRepository
func NewRepository(db *gorm.DB) Repository {
	repository := &GormAuditLogRepository{
		db: db,
	}
	return repository
}

// GormAuditLogRepository implements Repository using gorm
type GormAuditLogRepository struct {
	db *gorm.DB
}

// Create stores the given auditLog
func (r *GormAuditLogRepository) Create(ctx context.Context, auditLog *AuditLog) error {
	defer goa.MeasureSince([]string{"goa", "db", "auditLog", "create"}, time.Now())
	// check values
	if auditLog == nil {
		return errors.NewBadParameterErrorFromString("missing audit log auditLog to persist")
	}
	if auditLog.EventTypeID == uuid.Nil {
		return errors.NewBadParameterError("event_type_id", auditLog.EventTypeID)
	}
	if auditLog.IdentityID == uuid.Nil && auditLog.Username == "" {
		return errors.NewBadParameterErrorFromString("identity_id and username cannot be both missing at the same time")
	}
	db := r.db.Create(auditLog)
	if err := db.Error; err != nil {
		return errors.NewInternalError(ctx, err)
	}
	return nil
}

// LoadByID returns the AuditLog with the given id
// returns NotFoundError or InternalError
func (r *GormAuditLogRepository) LoadByID(ctx context.Context, id uuid.UUID) (AuditLog, error) {
	defer goa.MeasureSince([]string{"goa", "db", "auditLog", "loadById"}, time.Now())
	result := AuditLog{}
	tx := r.db.Model(AuditLog{}).Where("audit_log_id = ?", id).First(&result)
	if tx.RecordNotFound() {
		log.Error(nil, map[string]interface{}{
			"record_id": id,
		}, "auditlog not found")
		return AuditLog{}, errors.NewNotFoundError("auditlog_record", id.String())
	}
	if err := tx.Error; err != nil {
		return AuditLog{}, errors.NewInternalError(ctx, err)
	}
	return result, nil
}

// ListByIdentityID returns audit log records that belong to a user (given her identity ID), as well as the total number of records
// returns BadParameterError if the `start` or `limit` are invalid (negative) or InternalError an error if something wrong happened
// while qyerying or reading the returned rows
func (r *GormAuditLogRepository) ListByIdentityID(ctx context.Context, identityID uuid.UUID, start int, limit int) ([]AuditLog, int, error) {
	defer goa.MeasureSince([]string{"goa", "db", "auditLogs", "list_by_identity_id"}, time.Now())
	db := r.db.Model(&AuditLog{}).Where("identity_id = ?", identityID)
	unboundedDB := db
	if start < 0 {
		return nil, 0, errors.NewBadParameterError("start", start)
	}
	db = db.Offset(start)
	if limit <= 0 {
		return nil, 0, errors.NewBadParameterError("limit", limit)
	}
	db = db.Limit(limit)

	db = db.Select("count(*) over () as cnt2 , *")
	if err := db.Error; err != nil {
		return []AuditLog{}, 0, errors.NewInternalError(ctx, err)
	}
	rows, err := db.Rows()
	defer closeable.Close(ctx, rows)
	if err != nil {
		return []AuditLog{}, 0, errors.NewInternalError(ctx, err)
	}

	result := []AuditLog{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, errors.NewInternalError(ctx, err)
	}

	// need to set up a result for Scan() in order to extract total count.
	var count int
	var ignore interface{}
	columnValues := make([]interface{}, len(columns))

	for index := range columnValues {
		columnValues[index] = &ignore
	}
	columnValues[0] = &count
	first := true

	for rows.Next() {
		value := AuditLog{}
		db.ScanRows(rows, &value)
		if first {
			first = false
			if err = rows.Scan(columnValues...); err != nil {
				return nil, 0, errors.NewInternalError(ctx, err)
			}
		}
		result = append(result, value)
	}
	if first {
		// means 0 rows were returned from the first query (maybe becaus of offset outside of total count),
		// need to do a count(*) to find out total
		orgDB := unboundedDB.Select("count(*)")
		rows2, err := orgDB.Rows()
		defer closeable.Close(ctx, rows2)
		if err != nil {
			return nil, 0, errors.NewInternalError(ctx, err)
		}
		rows2.Next() // count(*) will always return a row
		rows2.Scan(&count)
	}
	return result, count, nil
}

// ListByUsername returns audit log records that belong to a user (given her username), as well as the total number of records
// returns BadParameterError if the `start` or `limit` are invalid (negative) or InternalError an error if something wrong happened
// while qyerying or reading the returned rows
func (r *GormAuditLogRepository) ListByUsername(ctx context.Context, username string, start int, limit int) ([]AuditLog, int, error) {
	defer goa.MeasureSince([]string{"goa", "db", "auditLogs", "list_by_username"}, time.Now())
	db := r.db.Model(&AuditLog{}).Where("username = ?", username)
	unboundedDB := db
	if start < 0 {
		return nil, 0, errors.NewBadParameterError("start", start)
	}
	db = db.Offset(start)
	if limit <= 0 {
		return nil, 0, errors.NewBadParameterError("limit", limit)
	}
	db = db.Limit(limit)

	db = db.Select("count(*) over () as cnt2 , *")
	if err := db.Error; err != nil {
		return []AuditLog{}, 0, errors.NewInternalError(ctx, err)
	}
	rows, err := db.Rows()
	defer closeable.Close(ctx, rows)
	if err != nil {
		return []AuditLog{}, 0, errors.NewInternalError(ctx, err)
	}

	result := []AuditLog{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, errors.NewInternalError(ctx, err)
	}

	// need to set up a result for Scan() in order to extract total count.
	var count int
	var ignore interface{}
	columnValues := make([]interface{}, len(columns))

	for index := range columnValues {
		columnValues[index] = &ignore
	}
	columnValues[0] = &count
	first := true

	for rows.Next() {
		value := AuditLog{}
		db.ScanRows(rows, &value)
		if first {
			first = false
			if err = rows.Scan(columnValues...); err != nil {
				return nil, 0, errors.NewInternalError(ctx, err)
			}
		}
		result = append(result, value)
	}
	if first {
		// means 0 rows were returned from the first query (maybe becaus of offset outside of total count),
		// need to do a count(*) to find out total
		orgDB := unboundedDB.Select("count(*)")
		rows2, err := orgDB.Rows()
		defer closeable.Close(ctx, rows2)
		if err != nil {
			return nil, 0, errors.NewInternalError(ctx, err)
		}
		rows2.Next() // count(*) will always return a row
		rows2.Scan(&count)
	}
	return result, count, nil
}
