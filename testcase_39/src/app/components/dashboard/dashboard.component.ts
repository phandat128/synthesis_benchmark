import { Component, OnInit } from '@angular/core';
import { SystemService } from '../../services/system.service';
import { AuthService } from '../../services/auth.service';

interface SystemStatus {
  status: 'OK' | 'DEGRADED' | 'CRITICAL';
  timestamp: string;
  details: string;
}

@Component({
  selector: 'app-dashboard',
  template: `
    <h2>System Health Dashboard</h2>

    <div *ngIf="status">
      <p>Current Status: <strong>{{ status.status }}</strong></p>
      <p>Last Checked: {{ status.timestamp | date:'medium' }}</p>
      <!-- Output Encoding is handled by Angular's interpolation by default -->
      <p>Details: {{ status.details }}</p>
    </div>
    <div *ngIf="!status">
      <p>Loading system status...</p>
    </div>

    <button (click)="checkStatus()">Refresh Status</button>

    <hr>

    <h3>Administrative Maintenance</h3>
    <p style="color: red;">Warning: System reset is a critical operation. Requires elevated privileges.</p>
    
    <!-- 
      SECURITY: The critical action is tied to a dedicated, authenticated service method. 
      It requires explicit user confirmation to prevent accidental or unauthorized execution.
    -->
    <button (click)="initiateSystemReset()" [disabled]="isResetting">
      {{ isResetting ? 'Resetting...' : 'Perform System Reset' }}
    </button>

    <div *ngIf="maintenanceMessage" [style.color]="maintenanceSuccess ? 'green' : 'red'">
      {{ maintenanceMessage }}
    </div>

    <button (click)="logout()">Logout</button>
  `,
  styles: ['button { margin: 5px; padding: 10px; }']
})
export class DashboardComponent implements OnInit {
  status: SystemStatus | null = null;
  isResetting: boolean = false;
  maintenanceMessage: string | null = null;
  maintenanceSuccess: boolean = false;

  constructor(
    private systemService: SystemService,
    private authService: AuthService
  ) { }

  ngOnInit(): void {
    this.checkStatus();
  }

  checkStatus(): void {
    this.systemService.checkSystemStatus().subscribe({
      next: (data) => {
        this.status = data;
      },
      error: (err) => {
        console.error('Failed to load status:', err);
        // Display a safe, non-sensitive error status
        this.status = { status: 'CRITICAL', timestamp: new Date().toISOString(), details: 'Failed to connect to monitoring service.' };
      }
    });
  }

  initiateSystemReset(): void {
    // Client-side confirmation for critical action (Defense in Depth)
    if (!confirm('Are you absolutely sure you want to reset the system? This action cannot be undone.')) {
      return;
    }

    this.isResetting = true;
    this.maintenanceMessage = 'Attempting system reset...';
    this.maintenanceSuccess = false;

    // Call the dedicated, authenticated maintenance endpoint
    this.systemService.performSystemReset().subscribe({
      next: (response) => {
        this.maintenanceMessage = `Success: ${response.message}`;
        this.maintenanceSuccess = true;
        this.isResetting = false;
        this.checkStatus(); // Refresh status after reset
      },
      error: (err) => {
        // Use the sanitized error message provided by the service
        this.maintenanceMessage = err.message || 'System reset failed due to an unknown error.';
        this.maintenanceSuccess = false;
        this.isResetting = false;
      }
    });
  }

  logout(): void {
    this.authService.logout();
  }
}