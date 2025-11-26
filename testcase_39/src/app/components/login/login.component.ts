import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AuthService } from '../../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  template: `
    <form [formGroup]="loginForm" (ngSubmit)="onSubmit()">
      <h2>System Login</h2>
      <div *ngIf="errorMessage" style="color: red;">{{ errorMessage }}</div>
      <div>
        <label for="username">Username:</label>
        <input id="username" type="text" formControlName="username" required>
        <div *ngIf="loginForm.get('username')?.invalid && loginForm.get('username')?.touched">
          Username is required and must be at least 3 characters.
        </div>
      </div>
      <div>
        <label for="password">Password:</label>
        <input id="password" type="password" formControlName="password" required>
        <div *ngIf="loginForm.get('password')?.invalid && loginForm.get('password')?.touched">
          Password is required and must be at least 8 characters.
        </div>
      </div>
      <button type="submit" [disabled]="loginForm.invalid">Login</button>
    </form>
  `,
  styles: ['div { margin-bottom: 10px; }']
})
export class LoginComponent {
  loginForm: FormGroup;
  errorMessage: string | null = null;

  constructor(
    private fb: FormBuilder,
    private authService: AuthService,
    private router: Router
  ) {
    // Implement robust client-side validation for input sanitation and integrity
    this.loginForm = this.fb.group({
      username: ['', [Validators.required, Validators.minLength(3)]],
      password: ['', [Validators.required, Validators.minLength(8)]]
    });
  }

  onSubmit(): void {
    this.errorMessage = null;

    if (this.loginForm.invalid) {
      this.errorMessage = 'Please fill out the form correctly.';
      // Mark all fields as touched to display validation messages
      this.loginForm.markAllAsTouched();
      return;
    }

    // Retrieve sanitized values from the form controls
    const credentials = {
      // Ensure username is trimmed before sending
      username: this.loginForm.value.username.trim(),
      password: this.loginForm.value.password
    };

    this.authService.login(credentials).subscribe({
      next: () => {
        this.router.navigate(['/dashboard']);
      },
      error: (err) => {
        // Display generic error message to prevent leaking backend details (Proper Error Handling)
        this.errorMessage = 'Login failed. Please check your credentials and try again.';
        console.error('Login attempt failed:', err);
      }
    });
  }
}