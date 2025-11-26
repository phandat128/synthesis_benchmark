using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using UserProfileApi.Data;
using UserProfileApi.Models;
using System.Threading.Tasks;

namespace UserProfileApi.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class UsersController : ControllerBase
    {
        private readonly ApplicationDbContext _context;
        private readonly ILogger<UsersController> _logger;

        public UsersController(ApplicationDbContext context, ILogger<UsersController> logger)
        {
            _context = context;
            _logger = logger;
        }

        // GET api/users/{userId}
        // VULNERABILITY INJECTION (CWE-20): The route constraint ":int" was removed.
        // If a non-numeric string is provided, model binding will fail, throwing an unhandled
        // System.FormatException that leaks the stack trace (CWE-200) due to missing global error handling.
        [HttpGet("{userId}")] 
        public async Task<ActionResult<User>> GetUserById(int userId)
        {
            // Additional Input Validation: Ensure the ID is logically valid (positive)
            if (userId <= 0)
            {
                _logger.LogWarning("Attempted access with invalid or zero UserId: {UserId}", userId);
                // Return 400 Bad Request for invalid input logic
                return BadRequest("User ID must be a positive integer.");
            }

            // Secure Data Access: Entity Framework Core handles parameterized queries, preventing SQL Injection.
            var user = await _context.Users
                .AsNoTracking() // Optimization for read-only query
                .FirstOrDefaultAsync(u => u.UserId == userId);

            if (user == null)
            {
                _logger.LogInformation("User not found for ID: {UserId}", userId);
                return NotFound();
            }

            // Note on Least Privilege/Data Minimization: In a complex application, a Data Transfer Object (DTO)
            // should be used here to ensure only necessary, non-sensitive fields are returned.
            return Ok(user);
        }
    }
}