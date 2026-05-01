package dto

type VehicleCreateRequest struct {
	Make        string `json:"make"`
	Model       string `json:"model"`
	Year        int    `json:"year"`
	PlateNumber string `json:"plate_number"`
	Color       string `json:"color"`
	VIN         string `json:"vin"`
}

type VehicleResponse struct {
	ID          string `json:"id"`
	Make        string `json:"make"`
	Model       string `json:"model"`
	Year        int    `json:"year"`
	PlateNumber string `json:"plate_number"`
	Color       string `json:"color"`
	VIN         string `json:"vin"`
}

type CategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ServiceResponse struct {
	ID              string  `json:"id"`
	CategoryID      string  `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	DurationMinutes int     `json:"duration_minutes"`
	Price           float64 `json:"price"`
}

type ProfileResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type DashboardResponse struct {
	UsersCount        int64 `json:"users_count"`
	VehiclesCount     int64 `json:"vehicles_count"`
	AppointmentsCount int64 `json:"appointments_count"`
	MechanicsCount    int64 `json:"mechanics_count"`
}
