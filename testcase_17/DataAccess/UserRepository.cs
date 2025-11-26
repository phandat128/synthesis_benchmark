using AuditSystem.Models;
using Microsoft.EntityFrameworkCore;
using System.Linq;
using System.Threading.Tasks;

namespace AuditSystem.DataAccess
{
    public interface IUserRepository
    {
        Task<User> AddUser(User user);
        Task<string?> GetUserNameById(int userId);
    }

    public class UserRepository : IUserRepository
    {
        // Using DbContext interface for flexibility, expecting AppDbContext implementation.
        private readonly DbContext _context;

        public UserRepository(DbContext context)
        {
            _context = context;
        }

        /// <summary>
        /// Safely adds a new user to the database.
        /// EF Core ensures parameterized queries, preventing First-Order SQL Injection.
        /// </summary>
        public async Task<User> AddUser(User user)
        {
            _context.Set<User>().Add(user);
            await _context.SaveChangesAsync();
            return user;
        }

        /// <summary>
        /// Retrieves the stored (potentially tainted) UserName based on UserId using safe LINQ.
        /// </summary>
        public async Task<string?> GetUserNameById(int userId)
        {
            return await _context.Set<User>()
                                 .Where(u => u.UserId == userId)
                                 .Select(u => u.UserName)
                                 .FirstOrDefaultAsync();
        }
    }
}