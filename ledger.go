package accountancy

import (
	"time"

	"github.com/go-xorm/xorm"
)

// Ledger -
type Ledger struct {
	ID      int64
	Dt      time.Time              `xorm:"notnull index"`
	DocID   int64                  `xorm:"notnull index"`
	DebID   int64                  `xorm:"notnull index"`
	CredID  int64                  `xorm:"notnull index"`
	Amount  int64                  `xorm:"notnull"`
	Props   map[string]interface{} `xorm:"jsonb"`
	Status  int                    `xorm:"notnull index"`
	Created time.Time              `xorm:"created"`
}

// MakeLedgerEntries -
func MakeLedgerEntries(eng *xorm.Engine, docID int64) (entries []*Ledger, err error) {

	return
}
