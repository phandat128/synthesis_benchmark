using Microsoft.EntityFrameworkCore;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;

namespace SecureDocManager.Data
{
    // Mock entity representing a sensitive report
    public class Report
    {
        public int Id { get; set; }
        public string Title { get; set; }
        public string Content { get; set; }
        public string RequiredGroups { get; set; } // Metadata for required groups
    }

    /// <summary>
    /// The Entity Framework Core context used to interface with the underlying user and report data store.
    /// </summary>
    public class ApplicationDbContext : DbContext
    {
        public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options)
            : base(options)
        {
        }

        public DbSet<Report> Reports { get; set; }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            base.OnModelCreating(modelBuilder);

            // Seed data for a sensitive report requiring both groups (A and B)
            modelBuilder.Entity<Report>().HasData(
                new Report { 
                    Id = 1, 
                    Title = "Q3 Financial Audit - Highly Restricted", 
                    Content = "Sensitive financial data only accessible by combined Group A and Group B members.",
                    RequiredGroups = "GroupClaimA,GroupClaimB"
                }
            );
        }
    }
}