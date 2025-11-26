import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn } from 'typeorm';

@Entity('users')
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  // The field where the potentially malicious input is stored.
  // TypeORM ensures safe persistence (1st-Order SQLi prevention).
  @Column({ unique: true, length: 50 })
  username: string;

  @CreateDateColumn()
  createdAt: Date;
}