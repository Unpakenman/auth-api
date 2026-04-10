{{define "CheckUserExist"}}
SELECT * from clinics.users
where phone = $1

{{end}}

{{define "GetEmployeeRoleByPhone"}}
SELECT r.name FROM clinics.role r
join clinics.employees e on r.role_id = e.role_id
where e.phone = $1
{{end}}

{{define "CreateUser"}}
INSERT INTO clinics.users (
  phone,
  email,
  password,
  is_employee,
  role
)VALUES(
  $1, $2, $3, $4,$5
)
{{end}}
