using Microsoft.AspNetCore.Components.Web;
using Microsoft.AspNetCore.Components.WebAssembly.Hosting;
using Microsoft.AspNetCore.Components.Authorization;
using Blazored.LocalStorage;
using DocumentManager.Providers;
using DocumentManager.Services;

namespace DocumentManager;

public class Program
{
    public static async Task Main(string[] args)
    {
        var builder = WebAssemblyHostBuilder.CreateDefault(args);
        builder.RootComponents.Add<App>("#app");
        builder.RootComponents.Add<HeadOutlet>("head::after");

        // 1. Configure HttpClient for API communication
        // Secure Coding: BaseAddress should point to the secure backend API endpoint.
        builder.Services.AddScoped(sp => new HttpClient { BaseAddress = new Uri(builder.HostEnvironment.BaseAddress) });

        // 2. Add Blazored Local Storage for secure token persistence
        builder.Services.AddBlazoredLocalStorage();

        // 3. Configure Custom Authentication Services
        builder.Services.AddAuthorizationCore();
        builder.Services.AddScoped<CustomAuthStateProvider>();
        // Register the custom provider as the official AuthenticationStateProvider
        builder.Services.AddScoped<AuthenticationStateProvider>(provider => provider.GetRequiredService<CustomAuthStateProvider>());
        
        // 4. Register Application Services
        builder.Services.AddScoped<AuthService>();

        await builder.Build().RunAsync();
    }
}