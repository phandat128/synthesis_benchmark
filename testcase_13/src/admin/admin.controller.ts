import { Controller, Get, Query, UsePipes, ValidationPipe, HttpStatus, HttpException } from '@nestjs/common';
import { AdminService } from './admin.service';
import { Activity } from '../activity/activity.entity';
import { IsString, Length, Matches } from 'class-validator';

// DTO for query parameter validation
class GetLogsDto {
  @IsString()
  @Length(3, 50)
  // Validate the input query parameter before passing it to the service.
  @Matches(/^[a-zA-Z0-9_]+$/, { message: 'Username query parameter is invalid.' })
  username: string;
}

@Controller('admin')
export class AdminController {
  constructor(private readonly adminService: AdminService) {}

  /**
   * Endpoint to retrieve activity logs for a specific user.
   * Input validation is performed on the query parameter.
   */
  @Get('activity-log')
  @UsePipes(new ValidationPipe({ transform: true, whitelist: true }))
  async getActivityLog(@Query() query: GetLogsDto): Promise<Activity[]> {
    try {
      // The validated username is passed to the service.
      return await this.adminService.fetchLogsByUsername(query.username);
    } catch (error) {
      if (error instanceof HttpException) {
        throw error;
      }
      console.error('Error fetching activity logs:', error);
      // SECURITY: Generic error message
      throw new HttpException('Failed to retrieve logs.', HttpStatus.INTERNAL_SERVER_ERROR);
    }
  }
}