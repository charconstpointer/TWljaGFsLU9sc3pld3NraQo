package repository

import (
	"fmt"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measurement"
)

func CreateMeasurement(db *[]*measurement.Measurement, m *measurement.Measurement) {
	*db = append(*db, m)
}

func GetMeasurements(db []*measurement.Measurement) []*measurement.Measurement {
	return db
}

func find(db []*measurement.Measurement, ID int) (*measurement.Measurement, error) {
	for _, m := range db {
		if m.ID == ID {
			return m, nil
		}
	}
	return nil, fmt.Errorf("requested measurement cannot be found %d", ID)
}
