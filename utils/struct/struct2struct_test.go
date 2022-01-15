package struct2struct

import (
	"fmt"
	"testing"
)

type User struct {
	Name string
	Role string
	//Age         int32
	EmployeCode int64 `copier:"EmployeNum"` // specify field name

	// Explicitly ignored in the destination struct.
	Salary int
}

type Employee struct {
	// Tell copier.Copy to panic if this field is not copied.
	Name string `copier:"must"`

	// Tell copier.Copy to return an error if this field is not copied.
	Age int32 `copier:"must,nopanic"`

	// Tell copier.Copy to explicitly ignore copying this field.
	Salary int `copier:"-"`

	DoubleAge int32
	EmployeId int64 `copier:"EmployeNum"` // specify field name
	SuperRole string
}

func TestStructCopy(t *testing.T) {
	var (
		user = User{Name: "Jinzhu", Role: "Admin", Salary: 200000}
		//user      = User{ Role: "Admin", Salary: 200000}
		//users     = []User{{Name: "Jinzhu", Age: 18, Role: "Admin", Salary: 100000}, {Name: "jinzhu 2", Age: 30, Role: "Dev", Salary: 60000}}
		employee = Employee{Salary: 150000}
		//employees = []Employee{}
	)
	err := StructCopy(&employee, &user)
	if err != nil {
		fmt.Println("ee ", err)
	}
	fmt.Printf("%#v \n", employee)
}
