import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { Router } from '@angular/router';

const TOKEN_KEY = 'authToken';
const API_URL = '/api/auth';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(private http: HttpClient, private router: Router) { }

  /**
   * Handles the login process and token storage.
   */
  login(credentials: any): Observable<{ token: string }> {
    // Basic validation check
    if (!credentials.username || !credentials.password) {
        return throwError(() => new Error('Username and password are required.'));
    }

    return this.http.post<{ token: string }>(`${API_URL}/login`, credentials)
      .pipe(
        tap(response => {
          this.storeToken(response.token);
        }),
        catchError(this.handleError)
      );
  }

  logout(): void {
    localStorage.removeItem(TOKEN_KEY);
    this.router.navigate(['/login']);
  }

  private storeToken(token: string): void {
    // Store token securely (using localStorage here, though HttpOnly cookies are often preferred for maximum security)
    const safeToken = token.trim();
    localStorage.setItem(TOKEN_KEY, safeToken);
  }

  getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY);
  }

  isLoggedIn(): boolean {
    const token = this.getToken();
    // Real implementation should include JWT expiry check
    return !!token;
  }

  private handleError(error: HttpErrorResponse) {
    let errorMessage = 'An unknown error occurred!';
    if (error.error instanceof ErrorEvent) {
      // Client-side network error
      errorMessage = `Client Error: ${error.error.message}`;
    } else {
      // Backend error. Avoid leaking internal details (Proper Error Handling).
      if (error.status === 401) {
        errorMessage = 'Authentication failed. Invalid credentials.';
      } else {
        errorMessage = `Server Error (${error.status}): Failed to connect or process request.`;
      }
    }
    console.error(errorMessage);
    // Return an observable with a generic, user-facing error message.
    return throwError(() => new Error('Login failed. Please check your credentials.'));
  }
}