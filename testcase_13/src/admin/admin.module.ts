import { Module, forwardRef } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { AdminController } from './admin.controller';
import { AdminService } from './admin.service';
import { Activity } from '../activity/activity.entity';
import { UserModule } from '../user/user.module';

@Module({
  imports: [
    TypeOrmModule.forFeature([Activity]),
    // Use forwardRef because AdminService is injected into UserService
    forwardRef(() => UserModule),
  ],
  controllers: [AdminController],
  providers: [AdminService],
  // Export AdminService so UserService can use it (circular dependency resolution)
  exports: [AdminService, TypeOrmModule],
})
export class AdminModule {}