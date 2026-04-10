package models

type UserExistResponse struct {
	UserID     int    `db:"user_id"`
	Phone      string `db:"phone"`
	Email      string `db:"email"`
	Password   string `db:"password"`
	IsEmployee bool   `db:"is_employee"`
	Role       string `db:"role"`
}

type EmployeeRoleResponse struct {
	Role string `db:"name"`
}
