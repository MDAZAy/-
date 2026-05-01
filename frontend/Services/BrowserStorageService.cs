using Microsoft.JSInterop;

namespace AutoService.Frontend.Services;

public sealed class BrowserStorageService(IJSRuntime jsRuntime)
{
    public ValueTask SetItemAsync(string key, string value) =>
        jsRuntime.InvokeVoidAsync("localStorage.setItem", key, value);

    public async Task<string?> GetItemAsync(string key) =>
        await jsRuntime.InvokeAsync<string?>("localStorage.getItem", key);

    public ValueTask RemoveItemAsync(string key) =>
        jsRuntime.InvokeVoidAsync("localStorage.removeItem", key);
}
