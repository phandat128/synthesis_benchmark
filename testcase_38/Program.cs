using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Authorization;
using Microsoft.IdentityModel.Tokens;
using SecureDocumentApi.Configuration;
using SecureDocumentApi.Filters;
using SecureDocumentApi.Services;
using System.Text;
using System.Threading.Tasks;

var builder = WebApplication.CreateBuilder(args);

// --- 1. Configure Services ---

// Register Authorization Service for DI (Contains the secure group check logic)
builder.Services.AddSingleton<AuthorizationService>();

builder.Services.AddControllers();

// Configure JWT Authentication
// NOTE: In a production environment, configuration values (Key, Issuer, Audience) must be loaded securely from secrets.
var jwtKey = builder.Configuration["Jwt:Key"] ?? "ThisIsAStrongSecretKeyForTestingPurposesOnly1234567890";
var issuer = builder.Configuration["Jwt:Issuer"] ?? "SecureDocumentIssuer";
var audience = builder.Configuration["Jwt:Audience"] ?? "SecureDocumentAudience";

builder.Services.AddAuthentication(options =>
{
    options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
    options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
})
.AddJwtBearer(options =>
{
    options.TokenValidationParameters = new TokenValidationParameters
    {
        ValidateIssuer = true,
        ValidateAudience = true,
        ValidateLifetime = true,
        ValidateIssuerSigningKey = true,
        ValidIssuer = issuer,
        ValidAudience = audience,
        IssuerSigningKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(jwtKey)),
        // Security Best Practice: Ensure RoleClaimType matches the claim used for groups/roles in the token.
        RoleClaimType = System.Security.Claims.ClaimTypes.Role 
    };
    
    // Proper Error Handling: Prevent sensitive information leakage on auth failure.
    options.Events = new JwtBearerEvents
    {
        OnAuthenticationFailed = context =>
        {
            context.Response.StatusCode = 401;
            context.Response.ContentType = "application/json";
            // Do not leak exception details in production
            var message = builder.Environment.IsDevelopment() ? context.Exception.Message : "Authentication failed due to invalid token.";
            return context.Response.WriteAsync($"{{ \"error\": \"{message}\" }}");
        }
    };
});

// Configure Authorization Policies
builder.Services.AddAuthorization(options =>
{
    // Define the strict policy requiring dual group membership (GROUP_A AND GROUP_B).
    options.AddPolicy(AuthorizationPolicies.DualGroupRequirement, policy =>
    {
        // This policy uses the custom requirement processed by RequiredGroupAuthorizationHandler.
        policy.AddRequirements(new RequiredGroupRequirement("GROUP_A", "GROUP_B"));
    });
});

// Register the custom Authorization Handler (the component that executes the secure group check logic).
builder.Services.AddSingleton<IAuthorizationHandler, RequiredGroupAuthorizationHandler>();

// Add Swagger/OpenAPI support
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();


// --- 2. Configure HTTP Request Pipeline ---

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}
else
{
    // Security Best Practice: Enforce HTTPS and HSTS in production
    app.UseHsts();
}

app.UseHttpsRedirection();

// 1. Authentication: Must run before Authorization
app.UseAuthentication();

// 2. Authorization: Checks permissions based on authenticated user
app.UseAuthorization();

app.MapControllers();

app.Run();