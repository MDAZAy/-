namespace AutoService.Frontend.Services;

public sealed class ApiEnvelope<T>
{
    public bool Success { get; set; }
    public T? Data { get; set; }
    public string Error { get; set; } = string.Empty;
    public string Code { get; set; } = string.Empty;
}

public sealed class AuthRequest
{
    public string Email { get; set; } = string.Empty;
    public string Password { get; set; } = string.Empty;
    public string Full_Name { get; set; } = string.Empty;
    public string Phone { get; set; } = string.Empty;
}

public sealed class AuthResponse
{
    public string Access_Token { get; set; } = string.Empty;
    public string Refresh_Token { get; set; } = string.Empty;
    public AuthUser User { get; set; } = new();
}

public sealed class AuthUser
{
    public string Id { get; set; } = string.Empty;
    public string Email { get; set; } = string.Empty;
    public string Full_Name { get; set; } = string.Empty;
    public string Phone { get; set; } = string.Empty;
    public string Role { get; set; } = string.Empty;
}

public sealed class ProfileResponse
{
    public string Id { get; set; } = string.Empty;
    public string Email { get; set; } = string.Empty;
    public string Full_Name { get; set; } = string.Empty;
    public string Phone { get; set; } = string.Empty;
    public string Role { get; set; } = string.Empty;
}

public sealed class VehicleRequest
{
    public string Make { get; set; } = string.Empty;
    public string Model { get; set; } = string.Empty;
    public int Year { get; set; } = DateTime.UtcNow.Year;
    public string Plate_Number { get; set; } = string.Empty;
    public string Color { get; set; } = string.Empty;
    public string Vin { get; set; } = string.Empty;
}

public sealed class VehicleResponse
{
    public string Id { get; set; } = string.Empty;
    public string Make { get; set; } = string.Empty;
    public string Model { get; set; } = string.Empty;
    public int Year { get; set; }
    public string Plate_Number { get; set; } = string.Empty;
    public string Color { get; set; } = string.Empty;
    public string Vin { get; set; } = string.Empty;
}

public sealed class CategoryResponse
{
    public string Id { get; set; } = string.Empty;
    public string Name { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
}

public sealed class ServiceResponse
{
    public string Id { get; set; } = string.Empty;
    public string Category_Id { get; set; } = string.Empty;
    public string Category_Name { get; set; } = string.Empty;
    public string Name { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public int Duration_Minutes { get; set; }
    public decimal Price { get; set; }
}

public sealed class AppointmentRequest
{
    public string Vehicle_Id { get; set; } = string.Empty;
    public string Service_Id { get; set; } = string.Empty;
    public string Start_Time { get; set; } = string.Empty;
    public string? End_Time { get; set; }
    public string Notes { get; set; } = string.Empty;
}

public sealed class AppointmentResponse
{
    public string Id { get; set; } = string.Empty;
    public string Confirmation_Number { get; set; } = string.Empty;
    public string Status { get; set; } = string.Empty;
    public ServiceResponse Service { get; set; } = new();
    public VehicleResponse Vehicle { get; set; } = new();
    public string Mechanic_Name { get; set; } = string.Empty;
    public string Start_Time_Utc { get; set; } = string.Empty;
    public string End_Time_Utc { get; set; } = string.Empty;
    public string Start_Time_Local { get; set; } = string.Empty;
    public string End_Time_Local { get; set; } = string.Empty;
    public string Notes { get; set; } = string.Empty;
}

public sealed class SlotResponse
{
    public string Start_Time_Utc { get; set; } = string.Empty;
    public string End_Time_Utc { get; set; } = string.Empty;
    public string Start_Time_Local { get; set; } = string.Empty;
    public string End_Time_Local { get; set; } = string.Empty;
}

public sealed class AvailableSlotsResponse
{
    public string Date { get; set; } = string.Empty;
    public string Timezone { get; set; } = string.Empty;
    public List<SlotResponse> Slots { get; set; } = [];
}

public sealed class DashboardResponse
{
    public long Users_Count { get; set; }
    public long Vehicles_Count { get; set; }
    public long Appointments_Count { get; set; }
    public long Mechanics_Count { get; set; }
}
