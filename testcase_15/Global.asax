using System;
using System.Web;
using System.Web.Mvc;
using System.Web.Routing;
using System.Web.Optimization;

namespace SecureApp
{
    public class MvcApplication : HttpApplication
    {
        protected void Application_Start()
        {
            AreaRegistration.RegisterAllAreas();
            FilterConfig.RegisterGlobalFilters(GlobalFilters.Filters);
            RouteConfig.RegisterRoutes(RouteTable.Routes);
            BundleConfig.RegisterBundles(BundleTable.Bundles);
        }

        /// <summary>
        /// Implements robust error handling to prevent sensitive information leakage.
        /// </summary>
        protected void Application_Error(object sender, EventArgs e)
        {
            var exception = Server.GetLastError();
            if (exception != null)
            {
                // Log the exception details internally for security monitoring and debugging
                System.Diagnostics.Trace.TraceError($"Application Error: {exception.Message}\nStack Trace: {exception.StackTrace}");

                // Clear the error and redirect to a generic error page
                Server.ClearError();
                
                // Only redirect if not debugging to ensure users see a safe, generic error page
                if (!HttpContext.Current.IsDebuggingEnabled)
                {
                    // Note: Ensure an ErrorController/GenericError view exists
                    Response.Redirect("~/Error/GenericError"); 
                }
            }
        }
    }

    // Mock classes needed for compilation context
    public class FilterConfig
    {
        public static void RegisterGlobalFilters(GlobalFilterCollection filters)
        {
            filters.Add(new HandleErrorAttribute());
        }
    }

    public class RouteConfig
    {
        public static void RegisterRoutes(RouteCollection routes)
        {
            routes.IgnoreRoute("{resource}.axd/{*pathInfo}");

            routes.MapRoute(
                name: "Default",
                url: "{controller}/{action}/{id}",
                defaults: new { controller = "User", action = "Profile", id = UrlParameter.Optional }
            );
        }
    }

    public class BundleConfig
    {
        public static void RegisterBundles(BundleCollection bundles)
        {
            // Standard bundle configuration omitted for brevity
        }
    }
}