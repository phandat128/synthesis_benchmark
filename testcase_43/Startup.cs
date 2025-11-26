using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using ReportApi.Data;
using ReportApi.Services;

namespace ReportApi
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
            // SECURE CODING: Configure API controllers.
            services.AddControllers();

            // Dependency Injection Setup
            // DataRepository is registered as Singleton as it primarily holds configuration (connection string)
            // and is thread-safe for executing commands.
            services.AddSingleton<DataRepository>();
            services.AddScoped<ReportGenerationService>();
            
            // Add logging
            services.AddLogging();
            
            // Note: In a real application, database context (if using EF Core) would be configured here.
            // services.AddDbContext<ApplicationDbContext>(options =>
            //     options.UseSqlServer(Configuration.GetConnectionString("DefaultConnection")));
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }
            else
            {
                // SECURE CODING: Use HSTS and HTTPS redirection in production.
                app.UseHsts();
                // app.UseExceptionHandler("/Error"); // Use a generic error handler
            }

            app.UseHttpsRedirection();
            app.UseRouting();

            // Authentication and Authorization middleware should be placed here if required.
            // app.UseAuthentication();
            // app.UseAuthorization();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
            });
        }
    }
}