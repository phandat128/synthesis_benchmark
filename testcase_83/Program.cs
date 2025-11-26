using Microsoft.AspNetCore.Components;
using Microsoft.AspNetCore.Components.Web;
using ResourceAllocationApp.Data;
using ResourceAllocationApp.Services;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddRazorPages();
builder.Services.AddServerSideBlazor();

// Register the secure Allocation Service for dependency injection.
// This service contains the critical logic for preventing Integer Overflow (CWE-190).
builder.Services.AddSingleton<AllocationService>();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (!app.Environment.IsDevelopment())
{
    // Use robust error handling without leaking sensitive details in production
    app.UseExceptionHandler("/Error");
    // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
    app.UseHsts();
}

app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();

app.MapBlazorHub();
app.MapFallbackToPage("/_Host");

app.Run();