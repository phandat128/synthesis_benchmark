using System.Web.Mvc;
using System.Web.Routing;
using System.Web;
using System;

namespace ProfileManager
{
    public class Global : HttpApplication
    {
        protected void Application_Start()
        {
            AreaRegistration.RegisterAllAreas();
            RouteConfig.RegisterRoutes(RouteTable.Routes);

            // Security Best Practice: Implement global error handling that prevents sensitive information leakage.
        }

        protected void Application_Error(object sender, EventArgs e)
        {
            var exception = Server.GetLastError();
            if (exception != null)
            {
                // Log the exception details securely (e.g., to a file or monitoring service)
                // System.Diagnostics.Trace.TraceError($"Global Error: {exception.Message}\n{exception.StackTrace}");

                // Clear the error and ensure a generic response
                Server.ClearError();
                
                // Prevent leaking stack traces or internal details to the client
                HttpContext.Current.Response.Clear();
                HttpContext.Current.Response.StatusCode = 500;
                
                // Set a generic error message or redirect to a friendly error page
                // HttpContext.Current.Response.Write("A critical error occurred.");
            }
        }
    }
}