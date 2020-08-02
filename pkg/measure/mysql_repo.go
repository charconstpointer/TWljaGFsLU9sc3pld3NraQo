package measure

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type MySQLRepo struct {
	DB *sqlx.DB
}

type entity struct {
	Id       int    `db:"Id"`
	Url      string `db:"Url"`
	Interval int    `db:"Delay"`
}

type probe struct {
	Id            int     `db:"Id"`
	Response      string  `db:"Response"`
	Duration      float32 `db:"Duration"`
	CreatedAt     int     `db:"CreatedAt"`
	MeasurementId int     `db:"MeasurementId"`
}

func (e entity) AsMeasure() *Measure {
	return &Measure{
		id:       e.Id,
		url:      e.Url,
		interval: e.Interval,
		probes:   make([]*Probe, 0),
	}
}

func (mr MySQLRepo) Save(m *Measure) (int, error) {
	q := "SELECT * FROM Measurements " +
		"WHERE Id = ?"
	var e []entity
	err := mr.DB.Select(&e, q, m.id)

	if len(e) == 0 {
		q = "INSERT INTO Measurements (Url, Delay)" +
			"VALUES (?,?)"

		res, err := mr.DB.Exec(q, m.url, m.interval)
		if err != nil {
			return -1, err
		}

		iid, err := res.LastInsertId()
		if err != nil {
			return -1, err
		}
		return int(iid), nil
	}
	ms := e[0]
	q = "UPDATE Measurements " +
		"SET Delay =? " +
		"WHERE Url =? "
	_, err = mr.DB.Exec(q, m.interval, m.url)
	if err != nil {
		return ms.Id, err
	}
	return ms.Id, nil
}

func (mr MySQLRepo) Get(ID int) (*Measure, error) {
	q := "SELECT * FROM Measurements " +
		"WHERE Id = ?"
	var e []entity
	err := mr.DB.Select(&e, q, ID)

	if len(e) == 0 {
		log.Err(err)
		return nil, err
	}
	m := e[0]
	q = "SELECT * FROM Probes " +
		"WHERE MeasurementId=?"
	var p []probe

	err = mr.DB.Select(&p, q, ID)
	measure := m.AsMeasure()

	for _, probe := range p {
		measure.AddProbe(NewProbe(probe.Response, probe.Duration, float32(probe.CreatedAt)))
	}

	return measure, nil
}

func (mr MySQLRepo) GetByUrl(URL string) (*Measure, error) {
	q := "SELECT * FROM Measurements " +
		"WHERE Measurements.Url=?"

	rows, err := mr.DB.Queryx(q, URL)

	if err != nil {
		log.Err(err)
		return nil, err
	}

	for rows.Next() {
		var p entity
		err = rows.StructScan(&p)
		if err != nil {
			log.Info().Msgf("%v\n", p)
			return nil, err
		}

		return p.AsMeasure(), nil

	}
	return nil, fmt.Errorf("could not query requested measure %s", URL)
}

func (mr MySQLRepo) GetAll() ([]*Measure, error) {
	var measures []entity
	q := "SELECT * FROM Measurements "
	err := mr.DB.Select(&measures, q)
	if err != nil {
		return nil, err
	}
	var m []*Measure
	for _, measure := range measures {
		m = append(m, measure.AsMeasure())
	}
	return m, nil
}

func (mr MySQLRepo) Update(ID int, interval int) error {
	q := "UPDATE Measurements " +
		"SET Delay = ? " +
		"WHERE Measurements.Id = ? "

	_, err := mr.DB.Exec(q, interval, ID)
	if err != nil {
		log.Err(err)
		return err
	}

	return nil
}

func (mr MySQLRepo) SaveProbe(ID int, p Probe) error {
	q := "INSERT INTO Probes (MeasurementId, Response, Duration, CreatedAt) " +
		"VALUES (?, ?, ?, ?)"

	_, err := mr.DB.Exec(q, ID, p.response, p.duration, p.createdAt)
	if err != nil {
		return err
	}
	return nil
}

func (mr MySQLRepo) Delete(ID int) error {
	q := "DELETE FROM Measurements " +
		"WHERE Measurements.Id = ? "

	res, err := mr.DB.Exec(q, ID)
	if err != nil {
		log.Err(err)
		return err
	}

	raf, err := res.RowsAffected()
	if int(raf) != 1 {
		return fmt.Errorf("something went wrong while deleting a measure %d", ID)
	}
	return nil
}
