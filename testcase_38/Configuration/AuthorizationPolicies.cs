namespace SecureDocumentApi.Configuration
{
    /// <summary>
    /// Defines named authorization policies used throughout the application.
    /// </summary>
    public static class AuthorizationPolicies
    {
        // Policy requiring membership in both Group A AND Group B for highly sensitive operations (e.g., DELETE).
        public const string DualGroupRequirement = "RequireDualGroupMembership";
    }
}