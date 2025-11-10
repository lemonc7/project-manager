package api

type EmployeeService interface {
}

type EmployeeHandler struct {
	employeeService EmployeeService
}
