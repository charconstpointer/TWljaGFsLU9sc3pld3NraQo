package repository

import (
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measurement"
)

func CreateMeasurement(db *[]*measurement.Measurement, m *measurement.Measurement) {
	*db = append(*db, m)
}

func GetMeasurements(db []*measurement.Measurement) []*measurement.Measurement {
	return db
}
