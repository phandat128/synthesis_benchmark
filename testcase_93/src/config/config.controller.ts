import { Controller, Post, Body, HttpStatus, HttpException, UsePipes, ValidationPipe } from '@nestjs/common';
import { ConfigService } from './config.service';
import { RuleExecutionDto } from './dto/rule-execution.dto';

@Controller('config')
export class ConfigController {
  constructor(private readonly configService: ConfigService) {}

  /**
   * Endpoint to securely execute a dynamic mathematical rule.
   * The input is automatically validated and sanitized by the global ValidationPipe
   * based on the constraints defined in RuleExecutionDto.
   */
  @Post('execute-rule')
  executeRule(@Body() ruleExecutionDto: RuleExecutionDto): { result: number } {
    try {
      // The service layer handles the secure execution logic.
      const result = this.configService.processRule(ruleExecutionDto);
      return { result };
    } catch (e) {
      // Proper Error Handling: Re-throw specific HTTP exceptions (like BadRequestException from the service)
      if (e instanceof HttpException) {
        throw e;
      }
      // Catch unexpected errors and return a generic 500 error to prevent leaking sensitive information.
      throw new HttpException('Internal server error during rule processing.', HttpStatus.INTERNAL_SERVER_ERROR);
    }
  }
}
