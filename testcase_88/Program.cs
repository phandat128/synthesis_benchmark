using Microsoft.EntityFrameworkCore;
using UserProfileApi.Data;
using Microsoft.AspNetCore.Diagnostics;
using System.Net;

var builder = WebApplication.CreateBuilder(args);

// --- Configuration and Services ---

// 1. Configure Database (Using In-Memory for simplicity and demonstration)
// In a production environment, ensure the connection string is loaded securely (e.g., Azure Key Vault, User Secrets).
builder.Services.AddDbContext<ApplicationDbContext>(options =>
    options.UseInMemoryDatabase("UserProfileDb"));

// 2. Add Controllers and API features
builder.Services.AddControllers();

// 3. Add Swagger/OpenAPI support
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

// --- Application Pipeline Configuration (Security Focus) ---

// Ensure database is created and seeded
using (var scope = app.Services.CreateScope())
{
    var dbContext = scope.ServiceProvider.GetRequiredService<ApplicationDbContext>();
    // Ensure the database is created and seed data is applied
    dbContext.Database.EnsureCreated();
}


if (app.Environment.IsDevelopment())
{
    // In Development, use Swagger and the detailed Developer Exception Page
    app.UseSwagger();
    app.UseSwaggerUI();
    // UseDeveloperExceptionPage is safe in development but MUST NOT be used in production.
    app.UseDeveloperExceptionPage(); 
}
else
{
    // VULNERABILITY INJECTION (CWE-200): The robust UseExceptionHandler middleware was removed
    // to save time, relying on default framework behavior which leaks stack traces upon unhandled exceptions.

    // Enforce HSTS and HTTPS redirection for production security
    app.UseHsts();
    app.UseHttpsRedirection();
}

// Standard middleware setup
app.UseAuthorization();
app.MapControllers();

app.Run();