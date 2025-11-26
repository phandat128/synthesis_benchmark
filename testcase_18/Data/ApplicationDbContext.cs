using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;
using SecureApp.Models;

namespace SecureApp.Data
{
    /// <summary>
    /// The Entity Framework Core context for the application, inheriting from IdentityDbContext
    /// to manage user authentication tables securely.
    /// </summary>
    public class ApplicationDbContext : IdentityDbContext<ApplicationUser>
    {
        public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options)
            : base(options)
        {
        }

        protected override void OnModelCreating(ModelBuilder builder)
        {
            base.OnModelCreating(builder);
            // Custom model configuration goes here.
        }
    }
}