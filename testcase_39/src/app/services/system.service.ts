import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';

const SYSTEM_API_URL = '/api/system';

interface SystemStatus {
  status: 'OK' | 'DEGRADED' | 'CRITICAL';
  timestamp: string;
  details: string;
}

@Injectable({
  providedIn: 'root'
})
export class SystemService {

  constructor(private http: HttpClient) { }

  /**
   * SECURE IMPLEMENTATION: Performs a read-only health check.
   * This method uses a simple GET request and explicitly avoids passing any query parameters
   * (like '?reset=true') that could trigger unauthorized state changes on the backend.
   * Since the AuthInterceptor requires a token for this path by default, this endpoint is protected.
   */
  checkSystemStatus(): Observable<SystemStatus> {
    return this.http.get<SystemStatus>(`${SYSTEM_API_URL}/health-check`)
      .pipe(
        catchError(this.handleError)
      );
  }

  /**
   * SECURE IMPLEMENTATION: Performs a critical system maintenance action (reset).
   * This uses a dedicated, authenticated, state-changing endpoint (POST) to ensure
   * the action is intentional and requires proper authorization (Least Privilege).
   */
  performSystemReset(): Observable<{ message: string }> {
    const maintenanceUrl = `${SYSTEM_API_URL}/maintenance/reset`;
    
    // Use POST for state-changing operations.
    return this.http.post<{ message: string }>(maintenanceUrl, { confirmation: true })
      .pipe(
        catchError(this.handleError)
      );
  }

  private handleError(error: HttpErrorResponse) {
    let errorMessage = 'System operation failed.';
    if (error.error instanceof ErrorEvent) {
      errorMessage = `Client Error: ${error.error.message}`;
    } else {
      // Server error. Do not leak internal stack traces.
      if (error.status === 403) {
        errorMessage = 'Authorization failed. You do not have permission for this maintenance action.';
      } else if (error.status === 401) {
        errorMessage = 'Authentication required.';
      } else {
        errorMessage = `Server Error (${error.status}): Failed to complete operation.`;
      }
    }
    console.error(errorMessage);
    return throwError(() => new Error(errorMessage));
  }
}