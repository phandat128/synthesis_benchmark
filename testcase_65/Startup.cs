using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using SecureApi.Services;
using Microsoft.AspNetCore.Http;
using System.Security.Claims;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Antiforgery;

namespace SecureApi
{
    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        // This method gets called by the runtime. Use this method to add services to the container.
        public void ConfigureServices(IServiceCollection services)
        {
            // Register services
            services.AddSingleton<IUserService, UserService>(); // Using Singleton for mock data

            // 1. Configure Authentication (using Cookies for session management)
            services.AddAuthentication("CookieAuth")
                .AddCookie("CookieAuth", options =>
                {
                    options.Cookie.Name = "SecureSession";
                    options.Cookie.HttpOnly = true; // Prevent client-side JS access to session cookie
                    options.Cookie.SecurePolicy = CookieSecurePolicy.Always; // Ensure cookies are only sent over HTTPS
                    // CRITICAL CSRF DEFENSE: SameSite=Strict prevents cookies from being sent with cross-site requests
                    options.Cookie.SameSite = SameSiteMode.Lax; 
                    options.ExpireTimeSpan = System.TimeSpan.FromMinutes(30);
                    options.SlidingExpiration = true;
                });

            // 2. Configure Anti-Forgery Services (Explicit defense against CSRF)
            services.AddAntiforgery(options =>
            {
                // The client will read the token from this cookie (set by the framework)
                options.Cookie.Name = "XSRF-REQUEST-TOKEN"; 
                options.Cookie.SameSite = SameSiteMode.Strict;
                options.Cookie.IsEssential = true;
                options.Cookie.HttpOnly = false; // Must be false so client-side JS can read the token value
                
                // The client must send the token back in this header name
                options.HeaderName = "X-CSRF-TOKEN"; 
            });

            services.AddControllers();
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env, IAntiforgery antiforgery)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }
            else
            {
                // Secure Error Handling: Use HSTS and redirect to non-leaking error page in production
                app.UseExceptionHandler("/Error");
                app.UseHsts();
            }

            app.UseHttpsRedirection();
            app.UseRouting();

            // Must be before UseAuthorization
            app.UseAuthentication(); 
            app.UseAuthorization();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
                
                // Mock Login Endpoint for testing cookie authentication
                endpoints.MapPost("/login", async context =>
                {
                    // Simulate successful authentication and set claims
                    var claims = new[]
                    {
                        new Claim(ClaimTypes.NameIdentifier, "1"), // User ID 1 (Alice)
                        new Claim(ClaimTypes.Name, "Alice")
                    };
                    var identity = new ClaimsIdentity(claims, "CookieAuth");
                    var principal = new ClaimsPrincipal(identity);

                    await context.SignInAsync("CookieAuth", principal);
                    context.Response.StatusCode = 200;
                    await context.Response.WriteAsync("Logged in successfully. Session cookie set.");
                });
            });
        }
    }
}