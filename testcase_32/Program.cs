using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using SecureReportApp.Services;
using SecureReportApp.Data;
using System.Security.Claims;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Logging;
using System.Threading.Tasks;

namespace SecureReportApp
{
    public class Program
    {
        public static void Main(string[] args)
        {
            var builder = WebApplication.CreateBuilder(args);

            // --- Security Configuration and Dependency Injection ---

            builder.Services.AddControllers();

            // Configure Mock Authentication Scheme for testing RBAC flow
            builder.Services.AddAuthentication("TestScheme")
                .AddScheme<AuthenticationSchemeOptions, TestAuthHandler>("TestScheme", options => { });
            
            // Configure Authorization
            builder.Services.AddAuthorization(options =>
            {
                // Define a default policy requiring authentication
                options.FallbackPolicy = options.DefaultPolicy;
            });

            // Dependency Injection for Services and Repositories
            builder.Services.AddSingleton<UserRepository>();
            builder.Services.AddScoped<AuthorizationService>();
            
            // --- Application Build and Pipeline ---

            var app = builder.Build();

            if (app.Environment.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }
            else
            {
                // Secure Error Handling: Prevent leaking sensitive stack traces in production
                app.UseExceptionHandler("/Error");
                app.UseHsts(); // Enforce HTTPS security
            }

            app.UseHttpsRedirection();
            app.UseRouting();

            // Authentication and Authorization must be called in this order
            app.UseAuthentication();
            app.UseAuthorization();

            app.MapControllers();

            // --- Mock Login Endpoints for Testing Authorization Scenarios ---
            
            // Scenario 1: User with BOTH required claims (Admin AND Finance) -> Should succeed (200)
            app.MapGet("/login/adminfinance", (HttpContext context) =>
            {
                var claims = new[]
                {
                    new Claim(ClaimTypes.NameIdentifier, "user123"),
                    new Claim(ClaimTypes.Name, "DualAccessUser"),
                    new Claim(ClaimTypes.Role, "Admin"), 
                    new Claim(ClaimTypes.Role, "Finance") 
                };
                var identity = new ClaimsIdentity(claims, "TestScheme");
                context.SignInAsync("TestScheme", new ClaimsPrincipal(identity));
                return Results.Ok("Logged in as DualAccessUser (Admin & Finance). Test: GET /Report/ViewSensitive/42");
            });
            
            // Scenario 2: User with ONLY ONE required claim (Admin only) -> Should fail (403)
            app.MapGet("/login/adminonly", (HttpContext context) =>
            {
                var claims = new[]
                {
                    new Claim(ClaimTypes.NameIdentifier, "user456"),
                    new Claim(ClaimTypes.Name, "AdminOnlyUser"),
                    new Claim(ClaimTypes.Role, "Admin")
                };
                var identity = new ClaimsIdentity(claims, "TestScheme");
                context.SignInAsync("TestScheme", new ClaimsPrincipal(identity));
                return Results.Ok("Logged in as AdminOnlyUser. Test: GET /Report/ViewSensitive/42 (Expected 403 Forbidden).");
            });

            app.Run();
        }
    }

    // Helper class for mock authentication handler to support the test scheme
    public class TestAuthHandler : AuthenticationHandler<AuthenticationSchemeOptions>
    {
        public TestAuthHandler(Microsoft.Extensions.Options.IOptionsMonitor<AuthenticationSchemeOptions> options, ILoggerFactory logger, System.Text.Encodings.Web.UrlEncoder encoder, ISystemClock clock)
            : base(options, logger, encoder, clock) { }

        protected override Task<AuthenticateResult> HandleAuthenticateAsync()
        {
            // If the user has already been signed in via the mock login endpoint, use that identity.
            if (Context.User.Identity?.IsAuthenticated == true)
            {
                var ticket = new AuthenticationTicket(Context.User, "TestScheme");
                return Task.FromResult(AuthenticateResult.Success(ticket));
            }

            // Otherwise, return no result (unauthenticated)
            return Task.FromResult(AuthenticateResult.NoResult());
        }
    }
}