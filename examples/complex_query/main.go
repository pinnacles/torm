package main

import (
	"database/sql"
	"log"
	"time"

	tsql "github.com/hatajoe/ttools/driver/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pinnacles/torm"

	_ "github.com/hatajoe/ttools/driver/sql/mysql"
)

var (
	db *sqlx.DB
)

type Org struct {
	ID            int64     `db:"id"`
	Name          string    `db:"name"`
	Address       string    `db:"address"`
	EstablishedAt time.Time `db:"established_at"`
}

func (t Org) TableName() string {
	return "orgs"
}

type User struct {
	ID    int64  `db:"id"`
	OrgID int64  `db:"org_id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

func (t User) TableName() string {
	return "users"
}

type OrgUser struct {
	OrgID    int64  `db:"org_id"`
	OrgName  string `db:"org_name"`
	UserID   int64  `db:"user_id"`
	UserName string `db:"user_name"`
}

func init() {
	tsql.Tracing(true)
	torm.Register(Org{})
	torm.Register(User{})
}

func main() {
	var err error
	tdb, err := sql.Open("ttools-mysql", "root:@/torm?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	db = sqlx.NewDb(tdb, "ttools-mysql")
	defer db.Close()

	builder := torm.NewBuilder(db)

	must(insertOrgsWithPreparedStatement(builder))
	must(bulkInsertUsers(builder))
	must(printOrgUsers(builder))
	must(deleteUsers(builder))
	must(deleteOrgs(builder))
	must(printOrgUsers(builder))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printOrgUsers(builder *torm.Builder) error {
	q := `
		SELECT
			orgs.id AS org_id,
			orgs.name AS org_name,
			users.id AS user_id,
			users.name AS user_name
		FROM
			users
		INNER JOIN orgs ON orgs.id = users.org_id
		WHERE
			users.age < :age
	`
	stmt, err := builder.Querier().PrepareNamed(q)
	if err != nil {
		return err
	}

	orgUsers := []OrgUser{}
	if err := stmt.Select(&orgUsers, torm.KV{"age": 30}); err != nil {
		return err
	}

	log.Println("####################")
	for _, u := range orgUsers {
		log.Printf("%#v\n", u)
	}
	log.Println("####################")
	log.Println()
	return nil
}

func insertOrgsWithPreparedStatement(builder *torm.Builder) error {
	established, err := time.Parse("_2 Jan 2006", "1 Jul 2019")
	if err != nil {
		return err
	}
	orgs := []Org{
		{
			Name:          "foo, Inc.",
			Address:       "Tokyo",
			EstablishedAt: established,
		},
		{
			Name:          "bar, Inc.",
			Address:       "Berin",
			EstablishedAt: established,
		},
		{
			Name:          "baz, Inc.",
			Address:       "SanFrancisco",
			EstablishedAt: established,
		},
	}
	stmt, err := builder.Querier().PrepareNamed(`INSERT INTO orgs(name, address, established_at) VALUES(:name, :address, :established_at)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, org := range orgs {
		_, err := stmt.Exec(org)
		if err != nil {
			return err
		}
	}
	return nil
}

func bulkInsertUsers(builder *torm.Builder) error {
	orgs := []Org{}
	if err := builder.Select().Query(&orgs); err != nil {
		return err
	}
	users := []User{
		{
			Name:  "foo",
			OrgID: orgs[0].ID,
			Email: "foo@example.com",
			Age:   10,
		},
		{
			Name:  "bar",
			OrgID: orgs[1].ID,
			Email: "bar@example.com",
			Age:   20,
		},
		{
			Name:  "baz",
			OrgID: orgs[2].ID,
			Email: "baz@example.com",
			Age:   30,
		},
	}
	_, err := builder.Querier().NamedExec(`INSERT INTO users(name, org_id, email, age) VALUES(:name, :org_id, :email, :age)`, users)
	if err != nil {
		return err
	}

	return nil
}

func deleteUsers(builder *torm.Builder) error {
	stmt, err := builder.Querier().PrepareNamed(`DELETE FROM users WHERE age < :age`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(torm.KV{"age": 50})
	if err != nil {
		return err
	}

	return nil
}

func deleteOrgs(builder *torm.Builder) error {
	stmt, err := builder.Querier().PrepareNamed(`DELETE FROM orgs WHERE id < :id`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(torm.KV{"id": 1000})
	if err != nil {
		return err
	}

	return nil
}
