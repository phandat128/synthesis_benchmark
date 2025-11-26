using System;
using System.Collections.Generic;

namespace Models
{
    /// <summary>
    /// Defines the structure for user session data.
    /// NOTE: The [Serializable] attribute is required for BinaryFormatter usage.
    /// </summary>
    [Serializable] // VULNERABILITY INJECTION: Required for BinaryFormatter
    public class UserSessionData
    {
        public string UserId { get; set; }
        public bool IsAuthenticated { get; set; }
        public DateTime LastActivity { get; set; }
        public List<string> Roles { get; set; } = new List<string>();
    }
}