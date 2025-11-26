using AuditSystem.DataAccess;
using AuditSystem.Models;
using AuditSystem.Services;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Configuration;

// --- Internal DbContext Definition (Required for EF Core Setup) ---
// In a production environment, this would be in DataAccess/AppDbContext.cs
public class AppDbContext : DbContext
{
    public AppDbContext(DbContextOptions<AppDbContext> options) : base(options) { }

    public DbSet<User> Users { get; set; }
    public DbSet<AuditLog> AuditLogs { get; set; }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        // Configure primary keys and indices
        modelBuilder.Entity<User>()
            .HasKey(u => u.UserId);

        modelBuilder.Entity<AuditLog>()
            .HasKey(l => l.LogId);

        // Indexing the UserName field for efficient lookups in audit logs
        modelBuilder.Entity<AuditLog>()
            .HasIndex(l => l.UserName);
    }
}
// -----------------------------------------------------------------


var builder = WebApplication.CreateBuilder(args);

// Add services to the container.

// 1. Configure Database Context using SQL Server
// NOTE: Connection string should be securely managed (e.g., Azure Key Vault, environment variables) in production.
var connectionString = builder.Configuration.GetConnectionString("DefaultConnection") 
                       ?? "Server=(localdb)\\mssqllocaldb;Database=SecureAuditDB;Trusted_Connection=True;MultipleActiveResultSets=true";

builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseSqlServer(connectionString));

// Register DbContext as the generic DbContext for repositories/services that need it
builder.Services.AddScoped<DbContext>(provider => provider.GetRequiredService<AppDbContext>());


// 2. Register Application Services and Repositories
builder.Services.AddScoped<IUserRepository, UserRepository>();
builder.Services.AddScoped<IAuditService, AuditService>();


// 3. Configure Controllers and API features
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(); // For testing endpoints

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();

    // Ensure database is created and migrated in development
    using (var scope = app.Services.CreateScope())
    {
        var dbContext = scope.ServiceProvider.GetRequiredService<AppDbContext>();
        // Use Migrate() for production, EnsureCreated() for simple testing
        dbContext.Database.EnsureCreated(); 
    }
}

// Standard security middleware setup
app.UseHttpsRedirection();
app.UseAuthorization();
app.MapControllers();

app.Run();