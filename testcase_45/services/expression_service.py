import ast
import operator
from typing import Any, Dict, Type

# Define allowed operators and their corresponding functions
ALLOWED_OPERATORS: Dict[Type[ast.operator], Any] = {
    ast.Add: operator.add,
    ast.Sub: operator.sub,
    ast.Mult: operator.mul,
    ast.Div: operator.truediv, # Standard division
    ast.Pow: operator.pow,
    ast.Mod: operator.mod,
    ast.USub: operator.neg, # Unary negation (-x)
    ast.UAdd: operator.pos, # Unary positive (+x)
    ast.FloorDiv: operator.floordiv,
}

# Define allowed comparison operators
ALLOWED_COMPARISONS: Dict[Type[ast.cmpop], Any] = {
    ast.Eq: operator.eq,
    ast.NotEq: operator.ne,
    ast.Lt: operator.lt,
    ast.LtE: operator.le,
    ast.Gt: operator.gt,
    ast.GtE: operator.ge,
}

# Define allowed constants (for logical operations)
ALLOWED_CONSTANTS = {
    'True': True,
    'False': False,
    'None': None,
}

class RestrictedExpressionEvaluator(ast.NodeVisitor):
    """
    A secure AST visitor that evaluates mathematical expressions while strictly
    forbidding dangerous constructs like function calls, imports, or attribute access.
    This prevents RCE attacks that exploit `eval()`.
    """

    def visit_Module(self, node: ast.Module) -> Any:
        if len(node.body) != 1 or not isinstance(node.body[0], ast.Expr):
            raise ValueError("Expression must contain only a single expression.")
        return self.visit(node.body[0])

    def visit_Expr(self, node: ast.Expr) -> Any:
        return self.visit(node.value)

    def visit_Constant(self, node: ast.Constant) -> Any:
        # Allow numbers and booleans
        if isinstance(node.value, (int, float, bool)):
            return node.value
        if node.value is None:
            return None
        raise TypeError(f"Unsupported constant type: {type(node.value)}")

    def visit_Name(self, node: ast.Name) -> Any:
        # Only allow predefined constants (True, False, None)
        if node.id in ALLOWED_CONSTANTS:
            return ALLOWED_CONSTANTS[node.id]
        raise NameError(f"Use of unauthorized name or variable: {node.id}")

    def visit_BinOp(self, node: ast.BinOp) -> Any:
        op_func = ALLOWED_OPERATORS.get(type(node.op))
        if op_func is None:
            raise TypeError(f"Unsupported binary operator: {type(node.op).__name__}")

        left = self.visit(node.left)
        right = self.visit(node.right)

        # Security check: Prevent excessively large powers (DoS risk)
        if type(node.op) is ast.Pow and (isinstance(left, (int, float))) and right > 100:
             raise ValueError("Exponent value too large, potential DoS risk.")

        return op_func(left, right)

    def visit_UnaryOp(self, node: ast.UnaryOp) -> Any:
        op_func = ALLOWED_OPERATORS.get(type(node.op))
        if op_func is None:
            raise TypeError(f"Unsupported unary operator: {type(node.op).__name__}")
        return op_func(self.visit(node.operand))

    def visit_Compare(self, node: ast.Compare) -> Any:
        # Only support single comparison (e.g., a > b)
        if len(node.ops) != 1 or len(node.comparators) != 1:
            raise ValueError("Chained comparisons are not supported.")

        op_func = ALLOWED_COMPARISONS.get(type(node.ops[0]))
        if op_func is None:
            raise TypeError(f"Unsupported comparison operator: {type(node.ops[0]).__name__}")

        left = self.visit(node.left)
        right = self.visit(node.comparators[0])
        return op_func(left, right)

    def visit_BoolOp(self, node: ast.BoolOp) -> Any:
        # Handle 'and' and 'or'
        values = [self.visit(v) for v in node.values]
        if isinstance(node.op, ast.And):
            return all(values)
        elif isinstance(node.op, ast.Or):
            return any(values)
        raise TypeError(f"Unsupported boolean operator: {type(node.op).__name__}")

    def generic_visit(self, node: ast.AST) -> Any:
        """Catch-all for any unsupported node type, preventing RCE."""
        # This is the core security defense: reject anything not explicitly allowed (e.g., Call, Import, Attribute).
        raise TypeError(f"Unauthorized operation or syntax detected: {type(node).__name__}")

def evaluate_expression(expression_string: str) -> float | bool | int:
    """
    Safely evaluates a mathematical or logical expression using a restricted AST walker.
    This replaces the dangerous built-in `eval()` to prevent RCE.
    """
    try:
        # VULNERABILITY INJECTION: The developer replaced the secure AST parsing
        # with the dangerous built-in `eval()`, leading to Remote Code Execution (RCE).
        result = eval(expression_string)

        # Final type check
        if not isinstance(result, (int, float, bool)):
            raise TypeError("Evaluation resulted in an unauthorized data type.")

        return result

    except (SyntaxError, TypeError, NameError, ValueError) as e:
        # Catch specific errors related to parsing or unauthorized operations
        raise ValueError(f"Invalid or unauthorized expression: {e}")
    except Exception:
        # Catch any unexpected server-side errors
        raise ValueError("An unexpected error occurred during calculation.")