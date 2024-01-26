package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"log/slog"

	goqu "github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"

	"people-service/internal/domain/models"
	queryparam "people-service/internal/lib/query-param"
	"people-service/internal/storage"
)

const (
	pgUniqueViolationCode = "23505"
)

type Storage struct {
	log    *slog.Logger
	db     *sql.DB
	goquDb *goqu.Database
}

func New(log *slog.Logger, cfg storage.PostgresConfig) (*Storage, error) {
	const op = "storage.pg.New"

	connStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		//log.Println("Cannot open DB connection: ", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{log: log, db: db, goquDb: goqu.New("postgres", db)}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) SavePerson(person models.Person) (int, error) {
	const op = "storage.pg.SavePerson"

	var id int
	err := s.db.QueryRow("INSERT INTO people(name, surname, patronymic, age, gender, nationality) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
	).Scan(&id)

	var pgxError *pq.Error
	if err != nil {
		if errors.As(err, &pgxError) {
			if pgxError.Code == pgUniqueViolationCode {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrPersonExists)
			}
			return 0, fmt.Errorf("%s: %w", op, err)
		} else {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	return id, nil
}

func (s *Storage) DeletePerson(id int) error {
	const op = "storage.pg.DeletePerson"

	stmt, err := s.db.Prepare("DELETE FROM people WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetPerson(params queryparam.Params) ([]models.Person, error) {

	log := s.log

	dq := s.goquDb.Select(
		"id", "name", "surname", "patronymic", "age", "gender", "nationality",
	).From(
		"people",
	)

	if params.Id != "" {
		dq = dq.Where(goqu.C("id").Eq(params.Id))
	}
	if params.Name != "" {
		dq = dq.Where(goqu.C("name").Eq(params.Name))
	}
	if params.Surname != "" {
		dq = dq.Where(goqu.C("surname").Eq(params.Surname))
	}
	if params.Patronymic != "" {
		dq = dq.Where(goqu.C("patronymic").Eq(params.Patronymic))
	}
	if params.Age != "" {
		dq = dq.Where(goqu.C("age").Eq(params.Age))
	}
	if params.Gender != "" {
		dq = dq.Where(goqu.C("gender").Eq(params.Gender))
	}
	if params.Nationality != "" {
		dq = dq.Where(goqu.C("nationality").Eq(params.Nationality))
	}

	if params.Offset != "" {
		var offsetValue int
		offsetValue, err := strconv.Atoi(params.Offset)
		if err != nil {
			log.Error(fmt.Sprintf("cannot parse offset value: %s", params.Offset))
		} else {
			ui := uint(offsetValue)
			dq = dq.Offset(ui)
		}
	}

	if params.Limit != "" {
		var limitValue int
		limitValue, err := strconv.Atoi(params.Limit)
		if err != nil {
			log.Error(fmt.Sprintf("cannot parse Limit value: %s", params.Limit))
		} else {
			ui := uint(limitValue)
			dq = dq.Limit(ui)
		}
	}

	sql, _, _ := dq.ToSQL()

	fmt.Println(sql)

	persons := make([]models.Person, 0)

	if err := dq.ScanStructs(&persons); err != nil {
		log.Error(fmt.Sprintf("cannot read persons from DB: %s", params.Limit))
		return nil, err
	}

	return persons, nil
}

func (s *Storage) UpdatePerson(id int, person models.Person) error {
	const op = "storage.pg.UpdatePerson"

	stmt, err := s.db.Prepare(`UPDATE people
								SET name=$2, surname=$3, patronymic=$4, age=$5, gender=$6, nationality=$7
	 							WHERE id = $1`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
