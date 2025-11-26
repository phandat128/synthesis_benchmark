using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using ImageProcessorApp.Models;
using ImageProcessorApp.Services;
using System;

namespace ImageProcessorApp.Controllers
{
    public class ImageController : Controller
    {
        private readonly ILogger<ImageController> _logger;
        private readonly ImageProcessorService _imageProcessorService;

        public ImageController(ILogger<ImageController> logger, ImageProcessorService imageProcessorService)
        {
            _logger = logger;
            _imageProcessorService = imageProcessorService;
        }

        [HttpGet]
        public IActionResult Upload()
        {
            return View();
        }

        [HttpPost]
        [ValidateAntiForgeryToken] // Standard security practice: Prevents CSRF
        public IActionResult ProcessImage(ProcessingRequest request)
        {
            // Input Validation: Check model state first (ensures required fields are present and types match)
            if (!ModelState.IsValid)
            {
                _logger.LogWarning("Invalid model state received for image processing.");
                return View("Upload", request);
            }

            // Input Validation: Ensure dimensions are positive and reasonable
            if (request.Width <= 0 || request.Height <= 0)
            {
                ModelState.AddModelError("", "Width and Height must be positive integers.");
                return View("Upload", request);
            }

            try
            {
                // Pass validated input to the service layer for secure size calculation
                bool success = _imageProcessorService.ProcessImage(request.Width, request.Height);

                if (success)
                {
                    TempData["Message"] = $"Image processing successful for {request.Width}x{request.Height}.";
                    return RedirectToAction("Upload");
                }
                else
                {
                    // This branch is hit if the service detects excessive size or resource limits
                    ModelState.AddModelError("", "Processing failed. Dimensions might be too large or invalid.");
                    return View("Upload", request);
                }
            }
            catch (ApplicationException ex)
            {
                // Catch the sanitized exception thrown by the service (robust error handling)
                _logger.LogError("Processing failed: {Message}", ex.Message);
                ModelState.AddModelError("", "An internal error prevented image processing. Please try smaller dimensions.");
                return View("Upload", request);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Unhandled exception during image processing.");
                // Generic error message to prevent information leakage
                ModelState.AddModelError("", "An unexpected error occurred.");
                return View("Upload", request);
            }
        }
    }
}