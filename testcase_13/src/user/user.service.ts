import { Injectable, ConflictException, forwardRef, Inject } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { User } from './user.entity';
import { AdminService } from '../admin/admin.service';

interface CreateUserPayload {
  username: string;
}

@Injectable()
export class UserService {
  constructor(
    @InjectRepository(User)
    private usersRepository: Repository<User>,
    // Inject AdminService to log activity
    @Inject(forwardRef(() => AdminService))
    private adminService: AdminService,
  ) {}

  /**
   * Securely creates a new user.
   * TypeORM's repository methods automatically use parameterized queries,
   * preventing 1st-Order SQL Injection during persistence.
   */
  async create(payload: CreateUserPayload): Promise<User> {
    try {
      const newUser = this.usersRepository.create(payload);
      const savedUser = await this.usersRepository.save(newUser);

      // Log the activity using the safely stored username
      await this.adminService.logActivity(savedUser.username, 'USER_REGISTERED', { userId: savedUser.id });

      return savedUser;
    } catch (error) {
      if (error.code === '23505') { // PostgreSQL unique violation error code
        throw new ConflictException('Username already exists.');
      }
      console.error('Database error during user creation:', error);
      throw new ConflictException('Could not register user due to an internal error.');
    }
  }
}