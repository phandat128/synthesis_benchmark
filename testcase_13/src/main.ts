import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ValidationPipe } from '@nestjs/common';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  // SECURITY BEST PRACTICE: Global Input Validation and Sanitization
  // Ensures all incoming request bodies and query parameters are validated
  // against defined DTOs, stripping unexpected fields (whitelist) and transforming types.
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true, // Strip properties not defined in the DTO
      forbidNonWhitelisted: true, // Throw error if extra properties are sent
      transform: true, // Automatically convert payload types
    }),
  );

  const port = 3000;
  await app.listen(port);
  console.log(`Application is running on: http://localhost:${port}`);
}
bootstrap();