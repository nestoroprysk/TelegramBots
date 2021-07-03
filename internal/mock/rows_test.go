package mock_test

import (
	"database/sql"
	"fmt"

	"github.com/nestoroprysk/TelegramBots/internal/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Person struct {
	Name string
	Age  int
}

var _ = It("Rows scan just fine", func() {
	r := mock.NewRows(mock.Cols("name", "age"),
		mock.Row("Peter", 30),
		mock.Row("George", 20),
	)

	var people []Person
	for r.Next() {
		p := Person{}
		Expect(r.Scan(&p.Name, &p.Age)).To(Succeed())
		people = append(people, p)
	}

	Expect(people).To(ConsistOf(
		Person{Name: "Peter", Age: 30},
		Person{Name: "George", Age: 20},
	))
})

var _ = It("Rows return columns", func() {
	r := mock.NewRows(mock.Cols("name", "age"))
	result, err := r.Columns()
	Expect(err).NotTo(HaveOccurred())
	Expect(result).To(ConsistOf("name", "age"))
})

var _ = It("Rows error is empty", func() {
	r := mock.NewRows()
	_, err := r.Columns()
	Expect(err).To(MatchError("no columns"))
})

var _ = It("Scan errors if nothing left", func() {
	r := mock.NewRows()
	err := r.Scan()
	Expect(err).To(MatchError(sql.ErrNoRows))
})

var _ = It("Next returns false if nothing left", func() {
	r := mock.NewRows()
	Expect(r.Next()).To(BeFalse())
})

var _ = It("Scanning too many columns does everything it may", func() {
	r := mock.NewRows(mock.Row("Peter"))
	p := Person{}
	Expect(r.Scan(&p.Name, &p.Age)).To(Succeed())
	Expect(p).To(Equal(Person{Name: "Peter"}))
})

var _ = It("Scanning too few columns does everything it may", func() {
	r := mock.NewRows(mock.Row("Peter", 10))
	p := Person{}
	Expect(r.Scan(&p.Name)).To(Succeed())
	Expect(p).To(Equal(Person{Name: "Peter"}))
})

var _ = It("Scanning errors as expected", func() {
	err := fmt.Errorf("oh no!")
	r := mock.NewRows(mock.RowErr(err))
	err = r.Scan(nil)
	Expect(err).To(MatchError("oh no!"))
})
