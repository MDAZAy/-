namespace AutoService.Frontend.Services;

public sealed class SessionState(BrowserStorageService storage)
{
    public string AccessToken { get; private set; } = string.Empty;
    public string RefreshToken { get; private set; } = string.Empty;
    public AuthUser CurrentUser { get; private set; } = new();
    public bool IsInitialized { get; private set; }
    public bool IsAuthenticated => !string.IsNullOrWhiteSpace(AccessToken);
    public bool IsAdmin => string.Equals(CurrentUser.Role, "admin", StringComparison.OrdinalIgnoreCase);

    public async Task InitializeAsync()
    {
        if (IsInitialized)
        {
            return;
        }

        AccessToken = await storage.GetItemAsync("access_token") ?? string.Empty;
        RefreshToken = await storage.GetItemAsync("refresh_token") ?? string.Empty;
        CurrentUser = new AuthUser
        {
            Id = await storage.GetItemAsync("user_id") ?? string.Empty,
            Email = await storage.GetItemAsync("user_email") ?? string.Empty,
            FullName = await storage.GetItemAsync("user_full_name") ?? string.Empty,
            Phone = await storage.GetItemAsync("user_phone") ?? string.Empty,
            Role = await storage.GetItemAsync("user_role") ?? string.Empty
        };

        IsInitialized = true;
    }

    public async Task SetSessionAsync(AuthResponse response)
    {
        AccessToken = response.AccessToken;
        RefreshToken = response.RefreshToken;
        CurrentUser = response.User;
        IsInitialized = true;

        await storage.SetItemAsync("access_token", AccessToken);
        await storage.SetItemAsync("refresh_token", RefreshToken);
        await storage.SetItemAsync("user_id", CurrentUser.Id);
        await storage.SetItemAsync("user_email", CurrentUser.Email);
        await storage.SetItemAsync("user_full_name", CurrentUser.FullName);
        await storage.SetItemAsync("user_phone", CurrentUser.Phone);
        await storage.SetItemAsync("user_role", CurrentUser.Role);
    }

    public async Task ClearAsync()
    {
        AccessToken = string.Empty;
        RefreshToken = string.Empty;
        CurrentUser = new AuthUser();

        await storage.RemoveItemAsync("access_token");
        await storage.RemoveItemAsync("refresh_token");
        await storage.RemoveItemAsync("user_id");
        await storage.RemoveItemAsync("user_email");
        await storage.RemoveItemAsync("user_full_name");
        await storage.RemoveItemAsync("user_phone");
        await storage.RemoveItemAsync("user_role");
    }
}
