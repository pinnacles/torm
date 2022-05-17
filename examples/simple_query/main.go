package main

import (
	"database/sql"
	"log"

	tsql "github.com/hatajoe/ttools/driver/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pinnacles/torm"

	_ "github.com/hatajoe/ttools/driver/sql/mysql"
)

var (
	db *sqlx.DB
)

type User struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

func (t User) TableName() string {
	return "users"
}

func init() {
	tsql.Tracing(true)
	torm.Register(User{})
}

func main() {
	var err error
	tdb, err := sql.Open("ttools-mysql", "root:@/torm")
	if err != nil {
		log.Fatal(err)
	}
	db = sqlx.NewDb(tdb, "ttools-mysql")
	defer db.Close()

	builder := torm.NewBuilder(db)

	must(insertUser(builder))
	must(printUsers(builder))
	must(insertUsers(builder))
	must(printUsersWhereByAge(builder))
	must(updateUser(builder))
	must(printUsers(builder))
	must(updateUsers(builder))
	must(printUsers(builder))
	must(deleteUser(builder))
	must(printUsers(builder))
	must(deleteUsers(builder))
	must(printUsers(builder))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printUsers(builder *torm.Builder) error {
	users := []User{}
	if err := builder.Select().Query(&users); err != nil {
		return err
	}
	log.Println("####################")
	for _, u := range users {
		log.Printf("%#v\n", u)
	}
	log.Println("####################")
	log.Println()
	return nil
}

func printUsersWhereByAge(builder *torm.Builder) error {
	users := []User{}
	if err := builder.Select().Where("age IN (:age)", torm.KV{"age": []int{20, 30}}).Query(&users); err != nil {
		return err
	}
	log.Println("####################")
	for _, u := range users {
		log.Printf("%#v\n", u)
	}
	log.Println("####################")
	log.Println()
	return nil
}

func insertUser(builder *torm.Builder) error {
	user := User{
		Name:  "foo",
		Email: "foo@example.com",
		Age:   10,
	}
	_, err := builder.Insert("name", "email", "age").Exec(user)
	return err
}

func insertUsers(builder *torm.Builder) error {
	users := []User{
		{
			Name:  "bar",
			Email: "bar@example.com",
			Age:   20,
		},
		{
			Name:  "baz",
			Email: "baz@example.com",
			Age:   30,
		},
	}
	for _, user := range users {
		if _, err := builder.Insert("name", "email", "age").Exec(user); err != nil {
			return err
		}
	}
	return nil
}

func updateUser(builder *torm.Builder) error {
	user := User{}
	if err := builder.Select().Where("name=:name", torm.KV{"name": "foo"}).Query(&user); err != nil {
		return err
	}
	user.Email = "fooooooo@example.com"
	_, err := builder.Update("email").Where("id=:id").Exec(user)
	return err
}

func updateUsers(builder *torm.Builder) error {
	users := []User{}
	if err := builder.Select().Where("name LIKE 'b%' AND age<=:age", torm.KV{"age": 20}).Query(&users); err != nil {
		return err
	}
	for i, u := range users {
		users[i].Age = u.Age + 3
	}
	for _, u := range users {
		if _, err := builder.Update("age", "email").Where("id=:id").Exec(u); err != nil {
			return err
		}
	}
	return nil
}

func deleteUser(builder *torm.Builder) error {
	user := User{}
	if err := builder.Select().Where("name=:name", torm.KV{"name": "foo"}).Query(&user); err != nil {
		return err
	}
	_, err := builder.Delete().Where("id=:id").Exec(user)
	return err
}

func deleteUsers(builder *torm.Builder) error {
	users := []User{}
	if err := builder.Select().Query(&users); err != nil {
		return err
	}
	for _, u := range users {
		if _, err := builder.Delete().Where("id=:id").Exec(u); err != nil {
			return err
		}
	}
	return nil
}
