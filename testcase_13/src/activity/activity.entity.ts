import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn } from 'typeorm';

@Entity('activity_logs')
export class Activity {
  @PrimaryGeneratedColumn()
  id: number;

  // This column stores the username used for lookup in the AdminService.
  @Column({ length: 50 })
  userName: string;

  @Column()
  action: string;

  @Column({ type: 'jsonb', nullable: true })
  details: any;

  @CreateDateColumn()
  timestamp: Date;
}