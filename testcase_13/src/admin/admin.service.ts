import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository, Connection } from 'typeorm';
import { Activity } from '../activity/activity.entity';
import { User } from '../user/user.entity';

@Injectable()
export class AdminService {
  constructor(
    @InjectRepository(Activity)
    private activityRepository: Repository<Activity>,
    @InjectRepository(User)
    private userRepository: Repository<User>,
    private connection: Connection, // Injected for completeness, but not used for raw queries
  ) {}

  /**
   * VULNERABLE IMPLEMENTATION: Introduces Second-Order SQL Injection.
   *
   * The username used here is retrieved from the database (potentially tainted).
   * We are now concatenating this stored, tainted value directly into a raw SQL query.
   */
  async fetchLogsByUsername(username: string): Promise<Activity[]> {
    // 1. Retrieve the user record safely (optional, but good for validation)
    const user = await this.userRepository.findOne({ where: { username } });
    if (!user) {
      throw new NotFoundException(`User ${username} not found.`);
    }

    // 2. VULNERABLE Log Retrieval using raw SQL concatenation (Sink)
    // The stored user.username is used unsafely.
    const rawQuery = `SELECT * FROM activity_logs WHERE "userName" = '${user.username}' ORDER BY "timestamp" DESC`;
    
    // DANGEROUS: Executing raw, concatenated query string.
    const logs = await this.connection.query(rawQuery);

    return logs;
  }

  // Helper function to log activity
  async logActivity(username: string, action: string, details?: any): Promise<void> {
    const log = this.activityRepository.create({
      userName: username,
      action: action,
      details: details,
    });
    await this.activityRepository.save(log);
  }
}