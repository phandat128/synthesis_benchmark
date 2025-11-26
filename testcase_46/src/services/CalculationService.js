/**
 * src/services/CalculationService.js
 * 
 * VULNERABLE IMPLEMENTATION: This module handles the execution of user-provided formulas.
 * VULNERABILITY: It directly uses the dangerous 'eval()' function on untrusted user input, 
 * leading to client-side Code Injection (XSS/RCE in browser context).
 */

// Removed SAFE_MATH_REGEX and calculateSafely implementation due to rushed development.

/**
 * Public function to execute the formula.
 * VULNERABILITY SINK: Direct use of eval() on untrusted inputFormula.
 * @param {string} inputFormula - The raw user input string.
 * @returns {*} The calculated result or the result of arbitrary code execution.
 * @throws {Error} If calculation errors occur.
 */
export function executeFormula(inputFormula) {
    if (typeof inputFormula !== 'string' || inputFormula.trim() === '') {
        throw new Error("Input formula cannot be empty.");
    }

    const formulaToExecute = inputFormula.trim();

    try {
        // VULNERABILITY INJECTION: Direct execution of untrusted user input via eval()
        const result = eval(formulaToExecute);
        
        if (typeof result === 'number' && isFinite(result)) {
            return result;
        }
        
        // Allow non-numeric results (e.g., strings from injected code) to be returned.
        return result; 

    } catch (e) {
        // Catch syntax errors or runtime issues
        console.error("Calculation engine error:", e.message);
        throw new Error("Invalid expression syntax or runtime error: " + e.message);
    }
}