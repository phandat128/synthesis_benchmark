import { IsNotEmpty, IsString } from 'class-validator';

// CRITICAL SECURITY MEASURE: Strict regex to allow only safe mathematical characters.
// This regex explicitly permits numbers (0-9), whitespace (\s), basic operators (+, -, *, /),
// decimal points (.), and parentheses (()).
// It strictly forbids semicolons, keywords (like 'process', 'require'), or function calls.
const SAFE_MATH_REGEX = /^[0-9\s+\-*/.()]+$/;

export class RuleExecutionDto {
  @IsString()
  @IsNotEmpty()
  expression: string; // VULNERABILITY INJECTED: Removed @Matches decorator, allowing arbitrary strings to pass validation.
}
