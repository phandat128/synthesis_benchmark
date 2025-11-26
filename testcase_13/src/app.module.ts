import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { UserModule } from './user/user.module';
import { AdminModule } from './admin/admin.module';
import { User } from './user/user.entity';
import { Activity } from './activity/activity.entity';

@Module({
  imports: [
    // WARNING: Credentials should be loaded from secure environment variables (e.g., ConfigService) in production.
    TypeOrmModule.forRoot({
      type: 'postgres',
      host: 'localhost',
      port: 5432,
      username: 'postgres',
      password: 'securepassword',
      database: 'logging_db',
      entities: [User, Activity],
      synchronize: true, // Set to false in production
      logging: false,
    }),
    UserModule,
    AdminModule,
  ],
})
export class AppModule {}