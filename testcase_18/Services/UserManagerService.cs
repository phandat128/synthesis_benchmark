using Microsoft.EntityFrameworkCore;
using SecureApp.Data;
using SecureApp.Models;
using System.Threading.Tasks;

namespace SecureApp.Services
{
    /// <summary>
    /// Contains the core business logic for user entity management.
    /// This layer assumes authorization checks have been completed by the Controller/API layer.
    /// </summary>
    public class UserManagerService
    {
        private readonly ApplicationDbContext _dbContext;

        public UserManagerService(ApplicationDbContext dbContext)
        {
            _dbContext = dbContext;
        }

        /// <summary>
        /// Retrieves a user by ID.
        /// </summary>
        /// <param name="userId">The ID of the user (string, as used by Identity).</param>
        /// <returns>The ApplicationUser or null if not found.</returns>
        public async Task<ApplicationUser> GetUserByIdAsync(string userId)
        {
            // Use FindAsync for efficient primary key lookup
            return await _dbContext.Users.FindAsync(userId);
        }

        /// <summary>
        /// Deletes a user by ID.
        /// </summary>
        /// <param name="userId">The ID of the user to delete.</param>
        /// <returns>True if deletion was successful or user was not found, false otherwise (e.g., database error).</returns>
        public async Task<bool> DeleteUserByIdAsync(string userId)
        {
            // 1. Input Validation: Ensure the ID is not null or empty.
            if (string.IsNullOrEmpty(userId))
            {
                return false;
            }

            var user = await GetUserByIdAsync(userId);

            if (user == null)
            {
                // User not found, operation is technically complete (idempotent success)
                return true; 
            }

            try
            {
                // Use EF Core tracking to remove the entity
                _dbContext.Users.Remove(user);
                // Use SaveChangesAsync to prevent blocking the thread
                await _dbContext.SaveChangesAsync();
                return true;
            }
            catch (DbUpdateConcurrencyException ex)
            {
                // Log the exception (CWE-778: Insufficient Logging)
                System.Console.WriteLine($"Database concurrency error during deletion: {ex.Message}");
                // Fail securely
                return false;
            }
            catch (DbUpdateException ex)
            {
                // Log the exception
                System.Console.WriteLine($"Database update error during deletion: {ex.Message}");
                return false;
            }
        }
    }
}