using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using ImageProcessorApp.Services;

namespace ImageProcessorApp
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
            // Register the secure Image Processor Service for dependency injection
            services.AddScoped<ImageProcessorService>();

            // Configure MVC with controllers and views
            services.AddControllersWithViews();

            // Configure Anti-Forgery tokens
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
                // Secure Error Handling: Use generic error page in production
                app.UseExceptionHandler("/Home/Error");
                // Enforce HSTS (HTTP Strict Transport Security)
                app.UseHsts();
            }

            // Standard security middleware
            app.UseHttpsRedirection();
            app.UseStaticFiles();

            app.UseRouting();

            // Authorization middleware would go here

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllerRoute(
                    name: "default",
                    pattern: "{controller=Image}/{action=Upload}/{id?}");
            });
        }
    }
}