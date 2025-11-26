using Microsoft.EntityFrameworkCore;
using UserProfileApi.Models;

namespace UserProfileApi.Data
{
    public class ApplicationDbContext : DbContext
    {
        public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options)
            : base(options)
        {
        }

        public DbSet<User> Users { get; set; }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            base.OnModelCreating(modelBuilder);

            // Seed data for initial testing
            modelBuilder.Entity<User>().HasData(
                new User { UserId = 1, Username = "Alice", Email = "alice@example.com" },
                new User { UserId = 2, Username = "Bob", Email = "bob@example.com" },
                new User { UserId = 3, Username = "Charlie", Email = "charlie@example.com" }
            );
        }
    }
}