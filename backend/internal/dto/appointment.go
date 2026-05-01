package dto

type AppointmentCreateRequest struct {
	VehicleID string  `json:"vehicle_id"`
	ServiceID string  `json:"service_id"`
	StartTime string  `json:"start_time"`
	EndTime   *string `json:"end_time,omitempty"`
	Notes     string  `json:"notes"`
}

type AppointmentResponse struct {
	ID                 string          `json:"id"`
	ConfirmationNumber string          `json:"confirmation_number"`
	Status             string          `json:"status"`
	Service            ServiceResponse `json:"service"`
	Vehicle            VehicleResponse `json:"vehicle"`
	MechanicName       string          `json:"mechanic_name"`
	StartTimeUTC       string          `json:"start_time_utc"`
	EndTimeUTC         string          `json:"end_time_utc"`
	StartTimeLocal     string          `json:"start_time_local"`
	EndTimeLocal       string          `json:"end_time_local"`
	Notes              string          `json:"notes"`
}

type SlotResponse struct {
	StartTimeUTC   string `json:"start_time_utc"`
	EndTimeUTC     string `json:"end_time_utc"`
	StartTimeLocal string `json:"start_time_local"`
	EndTimeLocal   string `json:"end_time_local"`
}

type AvailableSlotsResponse struct {
	Date     string         `json:"date"`
	Timezone string         `json:"timezone"`
	Slots    []SlotResponse `json:"slots"`
}
