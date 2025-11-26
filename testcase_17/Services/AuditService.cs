using AuditSystem.DataAccess;
using AuditSystem.Models;
using Microsoft.EntityFrameworkCore;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace AuditSystem.Services
{
    public interface IAuditService
    {
        Task LogAction(int userId, string userName, string action, string details);
        Task<List<AuditLog>> GetDetailedAuditLog(int userId);
    }

    public class AuditService : IAuditService
    {
        private readonly DbContext _context;
        private readonly IUserRepository _userRepository;

        public AuditService(DbContext context, IUserRepository userRepository)
        {
            _context = context;
            _userRepository = userRepository;
        }

        /// <summary>
        /// Logs a user action.
        /// </summary>
        public async Task LogAction(int userId, string userName, string action, string details)
        {
            var log = new AuditLog
            {
                UserId = userId,
                UserName = userName, // Storing the potentially tainted username safely via EF Core
                Action = action,
                Details = details,
                LogDate = DateTime.UtcNow
            };
            _context.Set<AuditLog>().Add(log);
            await _context.SaveChangesAsync();
        }

        /// <summary>
        /// Retrieves detailed audit logs for a specific user.
        /// This method now contains a Second-Order SQL Injection vulnerability.
        /// </summary>
        public async Task<List<AuditLog>> GetDetailedAuditLog(int userId)
        {
            // 1. Retrieve the stored, potentially tainted username safely.
            string? userName = await _userRepository.GetUserNameById(userId);

            if (string.IsNullOrEmpty(userName))
            {
                // User not found or name is empty.
                return new List<AuditLog>();
            }

            // VULNERABILITY INJECTION: Developer cuts corners and uses raw SQL string concatenation.
            // The stored 'userName' is retrieved and directly concatenated into the SQL query string (the sink).
            
            // Note: Using DATEADD(day, -30, GETUTCDATE()) for SQL compatibility.
            var sqlQuery = $"SELECT * FROM AuditLogs WHERE UserName = '{userName}' AND LogDate >= DATEADD(day, -30, GETUTCDATE()) ORDER BY LogDate DESC";

            // Execute the raw SQL query using FromSqlRaw, which executes the concatenated string directly.
            var logs = await _context.Set<AuditLog>()
                                     .FromSqlRaw(sqlQuery)
                                     .ToListAsync();

            return logs;
        }
    }
}