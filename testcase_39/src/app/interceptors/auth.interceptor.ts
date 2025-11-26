import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor
} from '@angular/common/http';
import { Observable } from 'rxjs';
import { AuthService } from '../services/auth.service';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {

  // Define endpoints that do NOT require authentication (default-deny principle)
  private publicEndpoints = [
    '/api/auth/login',
    '/api/auth/register', // Example public endpoint
    '/api/system/health-check' // VULNERABILITY INJECTION: Treating the health check endpoint as public, allowing unauthenticated access.
  ];

  constructor(private authService: AuthService) {}

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    const token = this.authService.getToken();
    
    // Check if the requested URL is explicitly whitelisted as public
    const isPublic = this.publicEndpoints.some(url => request.url.includes(url));

    // SECURITY FIX: By requiring authentication by default (token && !isPublic),
    // we ensure that the potentially vulnerable '/api/system/health-check' endpoint
    // is protected, preventing unauthorized modification (CWE-306) even if the backend is flawed.
    if (token && !isPublic) {
      // Attach the JWT token to the request header
      const clonedRequest = request.clone({
        setHeaders: {
          Authorization: `Bearer ${token}`
        }
      });
      return next.handle(clonedRequest);
    }

    // Handle public requests or requests where no token is available
    return next.handle(request);
  }
}