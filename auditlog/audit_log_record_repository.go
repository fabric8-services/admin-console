package auditlog

import (
	"context"
	"time"

	"github.com/fabric8-services/fabric8-common/closeable"
	"github.com/fabric8-services/fabric8-common/errors"
	"github.com/fabric8-services/fabric8-common/log"
	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"

	uuid "github.com/satori/go.uuid"
)

// RecordRepository provides functions to create and view Records
type RecordRepository interface {
	Create(ctx context.Context, record *Record) error
	LoadByID(ctx context.Context, id uuid.UUID) (Record, error)
	List(ctx context.Context, identityID uuid.UUID, start int, limit int) ([]Record, int, error)
}

// NewRecordRepository creates a GormRecordRepository
func NewRecordRepository(db *gorm.DB) RecordRepository {
	repository := &GormRecordRepository{
		db: db,
	}
	return repository
}

// GormRecordRepository implements RecordRepository using gorm
type GormRecordRepository struct {
	db *gorm.DB
}

// Create stores the given record
func (r *GormRecordRepository) Create(ctx context.Context, record *Record) error {
	defer goa.MeasureSince([]string{"goa", "db", "record", "create"}, time.Now())
	// check values
	if record == nil {
		return errors.NewBadParameterErrorFromString("missing audit log record to persist")
	}
	if record.IdentityID == uuid.Nil {
		return errors.NewBadParameterError("identity_id", record.IdentityID)
	}
	db := r.db.Create(record)
	if err := db.Error; err != nil {
		return errors.NewInternalError(ctx, err)
	}
	return nil
}

// LoadByID returns the record for the given id
// returns NotFoundError or InternalError
func (r *GormRecordRepository) LoadByID(ctx context.Context, id uuid.UUID) (Record, error) {
	defer goa.MeasureSince([]string{"goa", "db", "record", "loadById"}, time.Now())
	result := Record{}
	tx := r.db.Model(Record{}).Where("id = ?", id).First(&result)
	if tx.RecordNotFound() {
		log.Error(nil, map[string]interface{}{
			"record_id": id,
		}, "auditlog record not found")
		return Record{}, errors.NewNotFoundError("auditlog_record", id.String())
	}
	if err := tx.Error; err != nil {
		return Record{}, errors.NewInternalError(ctx, err)
	}
	return result, nil
}

// List returns audit log records that belong to a given user
func (r *GormRecordRepository) List(ctx context.Context, identityID uuid.UUID, start int, limit int) ([]Record, int, error) {
	defer goa.MeasureSince([]string{"goa", "db", "records", "list"}, time.Now())
	db := r.db.Model(&Record{}).Where("identity_id = ?", identityID)
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
		return []Record{}, 0, errors.NewInternalError(ctx, err)
	}
	rows, err := db.Rows()
	defer closeable.Close(ctx, rows)
	if err != nil {
		return []Record{}, 0, errors.NewInternalError(ctx, err)
	}

	result := []Record{}
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
		value := Record{}
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
