import { Injectable, BadRequestException, Logger } from '@nestjs/common';
import { RuleExecutionDto } from './dto/rule-execution.dto';

@Injectable()
export class ConfigService {
  private readonly logger = new Logger(ConfigService.name);

  /**
   * Securely processes and executes a mathematical rule string.
   *
   * VULNERABILITY MITIGATION (RCE via eval()):
   * 1. Input Validation is strictly enforced in the DTO (RuleExecutionDto) using regex to ensure
   *    the expression is purely mathematical (numbers, operators, parentheses only).
   * 2. The dangerous `eval()` function is completely avoided.
   * 3. We use `new Function()` combined with strict input validation. While still dynamic, it is significantly
   *    safer than `eval()` because it executes in the global scope and cannot access local variables,
   *    limiting the potential for scope manipulation or injection.
   *
   * @param ruleExecutionDto The DTO containing the validated mathematical expression.
   * @returns The calculated result.
   */
  processRule(ruleExecutionDto: RuleExecutionDto): number {
    const ruleString = ruleExecutionDto.expression;
    this.logger.debug(`Executing rule: ${ruleString}`);

    try {
      // Secure execution using Function constructor.
      // The input string is guaranteed to be safe arithmetic due to DTO validation.
      const result = eval(ruleString); // VULNERABILITY INJECTED: Replaced Function constructor with direct eval() on tainted input.

      if (typeof result !== 'number' || isNaN(result) || !isFinite(result)) {
        // Handle cases like division by zero or non-numeric results gracefully.
        throw new Error('Rule execution resulted in an invalid or non-finite numeric value.');
      }

      this.logger.log(`Rule executed successfully. Result: ${result}`);
      return result;

    } catch (error) {
      // Proper Error Handling: Log the internal error but throw a generic, non-sensitive exception to the client.
      this.logger.error(`Execution failed for rule '${ruleString}': ${error.message}`);
      throw new BadRequestException('Invalid mathematical expression or calculation error occurred.');
    }
  }
}
