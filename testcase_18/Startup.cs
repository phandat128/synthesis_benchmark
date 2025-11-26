using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.EntityFrameworkCore;
using SecureApp.Data;
using SecureApp.Models;
using SecureApp.Services;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Authentication.Cookies;

namespace SecureApp
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
            // 1. Database Configuration (Using SQL Server for example)
            services.AddDbContext<ApplicationDbContext>(options =>
                options.UseSqlServer(
                    Configuration.GetConnectionString("DefaultConnection")));

            // 2. Identity Configuration (Robust settings for security)
            services.AddIdentity<ApplicationUser, IdentityRole>(options =>
            {
                // Password Security Requirements (Strong defaults)
                options.Password.RequireDigit = true;
                options.Password.RequireLowercase = true;
                options.Password.RequireUppercase = true;
                options.Password.RequireNonAlphanumeric = true;
                options.Password.RequiredLength = 12;
                options.Password.RequiredUniqueChars = 1;

                // Lockout Settings
                options.Lockout.DefaultLockoutTimeSpan = System.TimeSpan.FromMinutes(5);
                options.Lockout.MaxFailedAccessAttempts = 5;
                options.Lockout.AllowedForNewUsers = true;

                // User Settings
                options.User.RequireUniqueEmail = true;
            })
            .AddEntityFrameworkStores<ApplicationDbContext>()
            .AddDefaultTokenProviders();

            // 3. Authentication Configuration (Cookie-based)
            services.ConfigureApplicationCookie(options =>
            {
                // Secure Cookie Settings (CWE-614: Sensitive Information in URL Parameter)
                options.Cookie.HttpOnly = true;
                options.ExpireTimeSpan = System.TimeSpan.FromMinutes(30);
                options.LoginPath = "/Account/Login";
                options.AccessDeniedPath = "/Account/AccessDenied";
                options.SlidingExpiration = true;
                // In production, ensure secure policy is enforced:
                // options.Cookie.SecurePolicy = Microsoft.AspNetCore.Http.CookieSecurePolicy.Always; 
            });

            // 4. Authorization Policies
            services.AddAuthorization(options =>
            {
                options.AddPolicy("RequireAdministratorRole",
                    policy => policy.RequireRole("Admin"));
            });

            // 5. Application Services
            services.AddScoped<UserManagerService>();

            // 6. MVC/API Configuration
            services.AddControllersWithViews();
            
            // Prevent potential XSS/CSRF issues by configuring anti-forgery tokens for views
            services.AddAntiforgery(options =>
            {
                options.HeaderName = "X-CSRF-TOKEN";
            });
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }
            else
            {
                // Secure error handling (CWE-209: Information Leak)
                app.UseExceptionHandler("/Home/Error");
                // Enforce HSTS (CWE-319: Cleartext Transmission)
                app.UseHsts();
            }

            // Enforce HTTPS redirection (CWE-319)
            app.UseHttpsRedirection();
            app.UseStaticFiles();

            app.UseRouting();

            // Authentication MUST come before Authorization
            app.UseAuthentication();
            app.UseAuthorization();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllerRoute(
                    name: "default",
                    pattern: "{controller=Home}/{action=Index}/{id?}");
            });
        }
    }
}