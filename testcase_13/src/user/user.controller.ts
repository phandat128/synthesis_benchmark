import { Controller, Post, Body, UsePipes, ValidationPipe, HttpStatus, HttpException } from '@nestjs/common';
import { UserService } from './user.service';
import { IsString, Length, Matches } from 'class-validator';

// DTO for input validation (Source of Taint)
// Ensures strict validation of the username format.
class CreateUserDto {
  @IsString()
  @Length(3, 50)
  // Defense-in-depth: Restrict characters to prevent common injection attempts
  @Matches(/^[a-zA-Z0-9_]+$/, { message: 'Username must contain only letters, numbers, and underscores.' })
  username: string;
}

@Controller('user')
export class UserController {
  constructor(private readonly userService: UserService) {}

  @Post('register')
  @UsePipes(new ValidationPipe({ whitelist: true, forbidNonWhitelisted: true }))
  async registerUser(@Body() createUserDto: CreateUserDto) {
    try {
      // Input is validated and sanitized by the DTO and ValidationPipe
      const user = await this.userService.create(createUserDto);
      return { id: user.id, username: user.username, message: 'User registered successfully.' };
    } catch (error) {
      if (error instanceof HttpException) {
        throw error;
      }
      // SECURITY: Generic error message to prevent leaking internal details
      throw new HttpException('Registration failed due to an internal error.', HttpStatus.INTERNAL_SERVER_ERROR);
    }
  }
}