package measure

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"time"
)

type MySQLRepo struct {
	DB *sqlx.DB
}

type entity struct {
	Id            int       `db:"Id"`
	Url           string    `db:"Url"`
	Response      string    `db:"Response"`
	Duration      float32   `db:"Duration"`
	CreatedAt     time.Time `db:"CreatedAt"`
	MeasurementId int       `db:"MeasurementId"`
}

type probe struct {
	Id            int       `db:"Id"`
	Response      string    `db:"Response"`
	Duration      float32   `db:"Duration"`
	CreatedAt     time.Time `db:"CreatedAt"`
	MeasurementId int       `db:"MeasurementId"`
}

func (mr *MySQLRepo) Save(m *Measure) error {
	q := "SELECT 1 FROM Measurements " +
		"WHERE Measurements.Url=?"

	rows, err := mr.DB.Queryx(q, m.url)
	if rows == nil || !rows.Next() || err != nil {
		q = "INSERT INTO Measurements (Url, Delay)" +
			"VALUES (?,?)"
		_, err := mr.DB.Exec(q, m.url, m.interval)
		if err != nil {
			return err
		}
		return nil
	}
	q = "UPDATE Measurements " +
		"SET Delay =? " +
		"WHERE Url =? "
	_, err = mr.DB.Exec(q, m.interval, m.url)
	if err != nil {
		return err
	}
	return nil
}

func (mr *MySQLRepo) Get(ID int) (*Measure, error) {
	q := "SELECT * FROM Measurements " +
		"JOIN Probes P on Measurements.Id = P.MeasurementId " +
		"WHERE Measurements.Id=?"

	rows, err := mr.DB.Queryx(q, ID)

	if err != nil {
		log.Err(err)
		return nil, err
	}

	for rows.Next() {
		var p entity
		err = rows.StructScan(&p)
		log.Info().Msgf("%v\n", p)
	}
	return nil, err

}

func (mr *MySQLRepo) GetByUrl(URL string) (*Measure, error) {
	panic("implement me")
}

func (mr *MySQLRepo) GetAll() ([]*Measure, error) {
	panic("implement me")
}

func (mr *MySQLRepo) Update(ID int, interval int) error {
	panic("implement me")
}

func (mr *MySQLRepo) Delete(ID int) error {
	panic("implement me")
}
