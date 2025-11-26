using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.EntityFrameworkCore;
using SecureDocManager.Data;
using SecureDocManager.Services;
using SecureDocManager.Filters;
using System.Security.Claims;
using Microsoft.AspNetCore.Authentication.Cookies;
using System.Collections.Generic;
using Microsoft.AspNetCore.Http;

namespace SecureDocManager
{
    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        public void ConfigureServices(IServiceCollection services)
        {
            // 1. Database Context (using In-Memory for demonstration)
            services.AddDbContext<ApplicationDbContext>(options =>
                options.UseInMemoryDatabase("SecureDocDB"));

            // 2. Custom Services
            services.AddScoped<UserService>();
            
            // 3. Authentication Setup (Mocking a simple cookie scheme)
            services.AddAuthentication(CookieAuthenticationDefaults.AuthenticationScheme)
                .AddCookie(options =>
                {
                    options.LoginPath = "/Account/Login";
                    options.AccessDeniedPath = "/Account/AccessDenied";
                    options.Cookie.HttpOnly = true; // Security best practice
                    options.Cookie.SecurePolicy = CookieSecurePolicy.Always; // Require HTTPS
                });

            // 4. MVC Configuration and Custom Filter Registration
            services.AddControllersWithViews();
            
            // Register the custom filter for dependency injection
            services.AddScoped<GroupAuthorizationFilter>();
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env, ApplicationDbContext dbContext)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }
            else
            {
                // Proper Error Handling: Implement robust error handling that does not leak sensitive information
                app.UseExceptionHandler("/Home/Error");
                app.UseHsts();
            }

            app.UseHttpsRedirection();
            app.UseStaticFiles();

            app.UseRouting();

            // Authentication and Authorization middleware must be in this order
            app.UseAuthentication();
            app.UseAuthorization();

            // Ensure the database is initialized and seeded (for in-memory DB)
            dbContext.Database.EnsureCreated();
            
            // Mock Login Endpoint for testing authorization logic
            app.UseEndpoints(endpoints =>
            {
                endpoints.MapGet("/login-mock/{groupA}/{groupB}", async context =>
                {
                    // Input validation for mock endpoint
                    if (!bool.TryParse(context.Request.RouteValues["groupA"]?.ToString(), out bool groupA)) groupA = false;
                    if (!bool.TryParse(context.Request.RouteValues["groupB"]?.ToString(), out bool groupB)) groupB = false;
                    
                    var claims = new List<Claim>
                    {
                        new Claim(ClaimTypes.NameIdentifier, "testuser123"),
                        new Claim(ClaimTypes.Name, $"User_A{groupA}_B{groupB}"),
                    };

                    if (groupA)
                    {
                        claims.Add(new Claim("GroupClaimA", "true"));
                    }
                    if (groupB)
                    {
                        claims.Add(new Claim("GroupClaimB", "true"));
                    }

                    var identity = new ClaimsIdentity(claims, CookieAuthenticationDefaults.AuthenticationScheme);
                    var principal = new ClaimsPrincipal(identity);

                    await context.SignInAsync(CookieAuthenticationDefaults.AuthenticationScheme, principal);
                    await context.Response.WriteAsync($"Logged in. Group A: {groupA}, Group B: {groupB}. Try accessing /reports/1");
                });

                endpoints.MapControllerRoute(
                    name: "default",
                    pattern: "{controller=Home}/{action=Index}/{id?}");
            });
        }
    }
}