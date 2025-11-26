using AuditSystem.DataAccess;
using AuditSystem.Models;
using AuditSystem.Services;
using Microsoft.AspNetCore.Mvc;
using System.ComponentModel.DataAnnotations;
using System.Threading.Tasks;

namespace AuditSystem.Controllers
{
    // DTO for incoming registration data, includes validation attributes.
    public record UserRegistrationDto(
        [Required(ErrorMessage = "Username is required.")] 
        [StringLength(50, MinimumLength = 3, ErrorMessage = "Username must be between 3 and 50 characters.")]
        string UserName
    );

    [ApiController]
    [Route("api/[controller]")]
    public class UserController : ControllerBase
    {
        private readonly IUserRepository _userRepository;
        private readonly IAuditService _auditService;

        public UserController(IUserRepository userRepository, IAuditService auditService)
        {
            _userRepository = userRepository;
            _auditService = auditService;
        }

        /// <summary>
        /// Registers a new user. This endpoint receives the initial tainted input (UserName).
        /// </summary>
        [HttpPost("register")]
        public async Task<IActionResult> RegisterUser([FromBody] UserRegistrationDto dto)
        {
            // 1. Input Validation: Ensure data conforms to expected constraints.
            if (!ModelState.IsValid)
            {
                return BadRequest(ModelState);
            }

            // 2. Normalization/Sanitization
            var sanitizedUserName = dto.UserName.Trim();

            var user = new User { UserName = sanitizedUserName };

            // 3. Safe Storage: UserRepository uses EF Core, which automatically parameterizes the insert query.
            var newUser = await _userRepository.AddUser(user);

            // Log the action
            await _auditService.LogAction(newUser.UserId, newUser.UserName, "USER_REGISTERED", $"New user registered.");

            return CreatedAtAction(nameof(RegisterUser), new { id = newUser.UserId }, new { newUser.UserId, newUser.UserName });
        }
    }
}