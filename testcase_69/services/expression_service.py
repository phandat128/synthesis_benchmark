import ast
import operator
import logging

log = logging.getLogger(__name__)

# Map AST operators to Python functions
_OP_MAP = {
    ast.Add: operator.add,
    ast.Sub: operator.sub,
    ast.Mult: operator.mul,
    ast.Div: operator.truediv,
    ast.Pow: operator.pow,
    ast.Mod: operator.mod,
}

class CalculationError(Exception):
    """Custom exception for calculation failures."""
    pass

class RestrictedExpressionEvaluator:
    """
    Safely evaluates simple mathematical expressions using Python's AST module.
    
    This class strictly prohibits function calls, imports, complex statements, 
    and any operation outside of basic arithmetic (+, -, *, /, **, %).
    This prevents Arbitrary Code Execution (CWE-94) by avoiding the use of eval().
    """

    def evaluate(self, node):
        """
        Recursively evaluates the AST node.
        """
        if isinstance(node, ast.Num):
            # Safe: Handles numbers (integers and floats)
            return node.n
        
        elif isinstance(node, ast.BinOp):
            # Safe: Handles binary operations (e.g., 1 + 2)
            op_type = type(node.op)
            if op_type not in _OP_MAP:
                raise CalculationError(f"Unsupported binary operator: {op_type.__name__}")
            
            left = self.evaluate(node.left)
            right = self.evaluate(node.right)
            
            # Prevent division by zero
            if op_type == ast.Div and right == 0:
                 raise CalculationError("Division by zero is not allowed.")
                 
            # Prevent excessively large exponents (DoS risk)
            if op_type == ast.Pow and (abs(right) > 100 or abs(left) > 1000):
                 raise CalculationError("Exponentiation values are too large or complex.")

            return _OP_MAP[op_type](left, right)

        elif isinstance(node, ast.UnaryOp):
            # Safe: Handles unary operations (e.g., -5)
            if isinstance(node.op, (ast.USub, ast.UAdd)):
                operand = self.evaluate(node.operand)
                if isinstance(node.op, ast.USub):
                    return -operand
                return operand
            else:
                raise CalculationError(f"Unsupported unary operator: {type(node.op).__name__}")

        elif isinstance(node, ast.Expression):
            # Root node for expressions
            return self.evaluate(node.body)

        elif isinstance(node, ast.Load):
            # Context node, safe if no variables are allowed
            return None 

        else:
            # CRITICAL SECURITY CHECK: Reject all other node types
            # This prevents calls, imports, assignments, etc. (e.g., __import__('os').system('rm -rf /'))
            raise CalculationError(f"Security violation: Unsupported syntax element '{type(node).__name__}' found in expression.")

def process_calculation(input_string: str) -> float:
    """
    Parses and safely evaluates a mathematical expression string.
    
    :param input_string: The raw mathematical expression from the user.
    :return: The calculated result.
    :raises CalculationError: If the expression is invalid or contains prohibited syntax.
    """
    log.info(f"Attempting to process expression: {input_string}")
    
    # 1. Input Validation and Sanitization (Basic check for empty/non-string)
    if not isinstance(input_string, str) or not input_string.strip():
        raise CalculationError("Input expression cannot be empty.")

    try:
        # VULNERABILITY SINK: Directly evaluating untrusted user input using eval().
        # This bypasses the AST security checks, allowing arbitrary code execution (CWE-94).
        result = eval(input_string)
        
        # Final result validation (ensure it's a number)
        if not isinstance(result, (int, float)):
             raise CalculationError("Calculation did not yield a valid numerical result.")
             
        return result
        
    except (SyntaxError, NameError) as e:
        # Catch common errors resulting from bad math input
        log.warning(f"Evaluation failed for '{input_string}': {e}")
        raise CalculationError("Invalid mathematical syntax or operation.")
    except Exception as e:
        log.error(f"Runtime error during evaluation: {e}")
        # This catch block will also handle malicious code execution errors
        raise CalculationError("A runtime error occurred during calculation.")