using System.Net.Http.Headers;
using System.Net.Http.Json;
using System.Text;
using System.Text.Json;

namespace AutoService.Frontend.Services;

public sealed class ApiClient(HttpClient httpClient, SessionState session)
{
    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNameCaseInsensitive = true
    };

    public async Task<AuthResponse> RegisterAsync(AuthRequest request) =>
        await SendAsync<AuthResponse>(HttpMethod.Post, "api/v1/auth/register", request);

    public async Task<AuthResponse> LoginAsync(AuthRequest request) =>
        await SendAsync<AuthResponse>(HttpMethod.Post, "api/v1/auth/login", request);

    public async Task<ProfileResponse> GetProfileAsync() =>
        await SendAsync<ProfileResponse>(HttpMethod.Get, "api/v1/me");

    public async Task<List<CategoryResponse>> GetCategoriesAsync() =>
        await SendAsync<List<CategoryResponse>>(HttpMethod.Get, "api/v1/service-categories");

    public async Task<List<ServiceResponse>> GetServicesAsync() =>
        await SendAsync<List<ServiceResponse>>(HttpMethod.Get, "api/v1/services");

    public async Task<List<VehicleResponse>> GetVehiclesAsync() =>
        await SendAsync<List<VehicleResponse>>(HttpMethod.Get, "api/v1/vehicles/my");

    public async Task<VehicleResponse> CreateVehicleAsync(VehicleRequest request) =>
        await SendAsync<VehicleResponse>(HttpMethod.Post, "api/v1/vehicles", request);

    public async Task<List<AppointmentResponse>> GetMyAppointmentsAsync() =>
        await SendAsync<List<AppointmentResponse>>(HttpMethod.Get, "api/v1/appointments/my");

    public async Task<List<AppointmentResponse>> GetAllAppointmentsAsync() =>
        await SendAsync<List<AppointmentResponse>>(HttpMethod.Get, "api/v1/appointments");

    public async Task<AvailableSlotsResponse> GetAvailableSlotsAsync(string date, string serviceId) =>
        await SendAsync<AvailableSlotsResponse>(HttpMethod.Get, $"api/v1/appointments/available-slots?date={Uri.EscapeDataString(date)}&service_id={Uri.EscapeDataString(serviceId)}");

    public async Task<AppointmentResponse> CreateAppointmentAsync(AppointmentRequest request) =>
        await SendAsync<AppointmentResponse>(HttpMethod.Post, "api/v1/appointments", request, Guid.NewGuid().ToString("N"));

    public async Task<DashboardResponse> GetDashboardAsync() =>
        await SendAsync<DashboardResponse>(HttpMethod.Get, "api/v1/admin/dashboard");

    private async Task<T> SendAsync<T>(HttpMethod method, string url, object? payload = null, string? idempotencyKey = null)
    {
        await session.InitializeAsync();

        using var request = new HttpRequestMessage(method, url);
        request.Headers.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
        if (session.IsAuthenticated)
        {
            request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", session.AccessToken);
        }
        if (!string.IsNullOrWhiteSpace(idempotencyKey))
        {
            request.Headers.Add("Idempotency-Key", idempotencyKey);
        }
        if (payload is not null)
        {
            request.Content = new StringContent(JsonSerializer.Serialize(payload), Encoding.UTF8, "application/json");
        }

        using var response = await httpClient.SendAsync(request);
        var raw = await response.Content.ReadAsStringAsync();
        var envelope = JsonSerializer.Deserialize<ApiEnvelope<T>>(raw, JsonOptions);

        if (!response.IsSuccessStatusCode || envelope is null || !envelope.Success || envelope.Data is null)
        {
            var error = envelope?.Error;
            if (string.IsNullOrWhiteSpace(error))
            {
                error = $"API error: {(int)response.StatusCode}";
            }
            throw new InvalidOperationException(error);
        }

        return envelope.Data;
    }
}
