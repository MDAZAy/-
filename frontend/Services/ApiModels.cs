using System.Text.Json.Serialization;

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
    [JsonPropertyName("email")]
    public string Email { get; set; } = string.Empty;

    [JsonPropertyName("password")]
    public string Password { get; set; } = string.Empty;

    [JsonPropertyName("full_name")]
    public string FullName { get; set; } = string.Empty;

    [JsonPropertyName("phone")]
    public string Phone { get; set; } = string.Empty;
}

public sealed class AuthResponse
{
    [JsonPropertyName("access_token")]
    public string AccessToken { get; set; } = string.Empty;

    [JsonPropertyName("refresh_token")]
    public string RefreshToken { get; set; } = string.Empty;

    [JsonPropertyName("user")]
    public AuthUser User { get; set; } = new();
}

public sealed class AuthUser
{
    [JsonPropertyName("id")]
    public string Id { get; set; } = string.Empty;

    [JsonPropertyName("email")]
    public string Email { get; set; } = string.Empty;

    [JsonPropertyName("full_name")]
    public string FullName { get; set; } = string.Empty;

    [JsonPropertyName("phone")]
    public string Phone { get; set; } = string.Empty;

    [JsonPropertyName("role")]
    public string Role { get; set; } = string.Empty;
}

public sealed class ProfileResponse
{
    [JsonPropertyName("id")]
    public string Id { get; set; } = string.Empty;

    [JsonPropertyName("email")]
    public string Email { get; set; } = string.Empty;

    [JsonPropertyName("full_name")]
    public string FullName { get; set; } = string.Empty;

    [JsonPropertyName("phone")]
    public string Phone { get; set; } = string.Empty;

    [JsonPropertyName("role")]
    public string Role { get; set; } = string.Empty;
}

public sealed class VehicleRequest
{
    [JsonPropertyName("make")]
    public string Make { get; set; } = string.Empty;

    [JsonPropertyName("model")]
    public string Model { get; set; } = string.Empty;

    [JsonPropertyName("year")]
    public int Year { get; set; } = DateTime.UtcNow.Year;

    [JsonPropertyName("plate_number")]
    public string PlateNumber { get; set; } = string.Empty;

    [JsonPropertyName("color")]
    public string Color { get; set; } = string.Empty;

    [JsonPropertyName("vin")]
    public string Vin { get; set; } = string.Empty;
}

public sealed class VehicleResponse
{
    [JsonPropertyName("id")]
    public string Id { get; set; } = string.Empty;

    [JsonPropertyName("make")]
    public string Make { get; set; } = string.Empty;

    [JsonPropertyName("model")]
    public string Model { get; set; } = string.Empty;

    [JsonPropertyName("year")]
    public int Year { get; set; }

    [JsonPropertyName("plate_number")]
    public string PlateNumber { get; set; } = string.Empty;

    [JsonPropertyName("color")]
    public string Color { get; set; } = string.Empty;

    [JsonPropertyName("vin")]
    public string Vin { get; set; } = string.Empty;
}

public sealed class CategoryResponse
{
    [JsonPropertyName("id")]
    public string Id { get; set; } = string.Empty;

    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;

    [JsonPropertyName("description")]
    public string Description { get; set; } = string.Empty;
}

public sealed class ServiceResponse
{
    [JsonPropertyName("id")]
    public string Id { get; set; } = string.Empty;

    [JsonPropertyName("category_id")]
    public string CategoryId { get; set; } = string.Empty;

    [JsonPropertyName("category_name")]
    public string CategoryName { get; set; } = string.Empty;

    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;

    [JsonPropertyName("description")]
    public string Description { get; set; } = string.Empty;

    [JsonPropertyName("duration_minutes")]
    public int DurationMinutes { get; set; }

    [JsonPropertyName("price")]
    public decimal Price { get; set; }
}

public sealed class AppointmentRequest
{

    [JsonPropertyName("vehicle_id")]
    public string VehicleId { get; set; } = string.Empty;

    [JsonPropertyName("service_id")]
    public string ServiceId { get; set; } = string.Empty;

    [JsonPropertyName("start_time")]
    public string StartTime { get; set; } = string.Empty;

    [JsonPropertyName("end_time")]
    public string? EndTime { get; set; }

    [JsonPropertyName("notes")]
    public string Notes { get; set; } = string.Empty;
}

public sealed class AppointmentResponse
{
    [JsonPropertyName("id")]
    public string Id { get; set; } = string.Empty;

    [JsonPropertyName("confirmation_number")]
    public string ConfirmationNumber { get; set; } = string.Empty;

    [JsonPropertyName("status")]
    public string Status { get; set; } = string.Empty;

    [JsonPropertyName("service")]
    public ServiceResponse Service { get; set; } = new();

    [JsonPropertyName("vehicle")]
    public VehicleResponse Vehicle { get; set; } = new();

    [JsonPropertyName("mechanic_name")]
    public string MechanicName { get; set; } = string.Empty;

    [JsonPropertyName("start_time_utc")]
    public string StartTimeUtc { get; set; } = string.Empty;

    [JsonPropertyName("end_time_utc")]
    public string EndTimeUtc { get; set; } = string.Empty;

    [JsonPropertyName("start_time_local")]
    public string StartTimeLocal { get; set; } = string.Empty;

    [JsonPropertyName("end_time_local")]
    public string EndTimeLocal { get; set; } = string.Empty;

    [JsonPropertyName("notes")]
    public string Notes { get; set; } = string.Empty;
}

public sealed class SlotResponse
{
    [JsonPropertyName("start_time_utc")]
    public string StartTimeUtc { get; set; } = string.Empty;

    [JsonPropertyName("end_time_utc")]
    public string EndTimeUtc { get; set; } = string.Empty;

    [JsonPropertyName("start_time_local")]
    public string StartTimeLocal { get; set; } = string.Empty;

    [JsonPropertyName("end_time_local")]
    public string EndTimeLocal { get; set; } = string.Empty;
}

public sealed class AvailableSlotsResponse
{
    [JsonPropertyName("date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("timezone")]
    public string Timezone { get; set; } = string.Empty;

    [JsonPropertyName("slots")]
    public List<SlotResponse> Slots { get; set; } = [];
}

public sealed class DashboardResponse
{
    [JsonPropertyName("users_count")]
    public long UsersCount { get; set; }

    [JsonPropertyName("vehicles_count")]
    public long VehiclesCount { get; set; }

    [JsonPropertyName("appointments_count")]
    public long AppointmentsCount { get; set; }

    [JsonPropertyName("mechanics_count")]
    public long MechanicsCount { get; set; }
}
