using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace SecureSessionApi
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
            // Configure API controllers. Use System.Text.Json as the default serializer.
            services.AddControllers()
                    .AddJsonOptions(options =>
                    {
                        // Enforce secure JSON serialization defaults
                        options.JsonSerializerOptions.PropertyNamingPolicy = System.Text.Json.JsonNamingPolicy.CamelCase;
                        options.JsonSerializerOptions.WriteIndented = false;
                    });

            // Configure secure data protection for cookies/session management (if used later)
            services.AddDataProtection();
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                // In production, detailed error pages should be disabled to prevent information leakage.
                app.UseDeveloperExceptionPage();
            }
            else
            {
                // Use HSTS and HTTPS redirection in production
                app.UseHsts();
            }

            app.UseHttpsRedirection();

            app.UseRouting();

            // Standard authorization/authentication middleware placement
            app.UseAuthorization();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
            });
        }
    }
}