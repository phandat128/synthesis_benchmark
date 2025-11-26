<template>
  <div class="formula-input bg-white shadow-lg rounded-lg p-6">
    <h2 class="text-2xl font-semibold mb-4">Enter Formula</h2>
    
    <div class="mb-4">
      <label for="formula" class="block text-sm font-medium text-gray-700 mb-1">Formula String (e.g., (10 + 5) * 2)</label>
      <textarea
        id="formula"
        v-model="formulaString"
        rows="4"
        class="w-full p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
        placeholder="Enter your mathematical expression here..."
      ></textarea>
    </div>

    <button 
      @click="calculateAndDisplay"
      :disabled="isCalculating"
      class="w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
    >
      {{ isCalculating ? 'Calculating...' : 'Calculate Result' }}
    </button>

    <!-- Display Area -->
    <div v-if="result !== null || error" class="mt-6 p-4 rounded-md" :class="{'bg-green-100 border border-green-400 text-green-700': result !== null, 'bg-red-100 border border-red-400 text-red-700': error}">
      <p v-if="result !== null" class="font-bold">Result:</p>
      <p v-if="result !== null" class="text-xl break-words">{{ result }}</p>
      
      <!-- Secure Error Handling: Displaying sanitized error message -->
      <p v-if="error" class="font-bold">Error:</p>
      <p v-if="error" class="break-words">{{ error }}</p>
    </div>

    <div class="mt-4 text-sm text-gray-500">
        <p>Note: This calculator only permits basic arithmetic operators (+, -, *, /) and numbers. Any attempt to inject code or use complex functions will be rejected.</p>
    </div>
  </div>
</template>

<script>
import { executeFormula } from '@/services/CalculationService';

export default {
  name: 'FormulaInput',
  data() {
    return {
      formulaString: '100 / (5 + 5) * 2', // Default example formula
      result: null,
      error: null,
      isCalculating: false,
    };
  },
  methods: {
    async calculateAndDisplay() {
      this.result = null;
      this.error = null;
      this.isCalculating = true;

      // Source of Taint: User input captured by v-model
      const inputFormula = this.formulaString;

      if (!inputFormula.trim()) {
        this.error = "Please enter a formula.";
        this.isCalculating = false;
        return;
      }

      try {
        // Propagation: Calling the service layer.
        // Mitigation occurs inside executeFormula via strict whitelisting and avoiding eval().
        const calculatedResult = executeFormula(inputFormula);
        
        // Output Encoding: Ensure the numeric result is treated as safe text for display.
        this.result = calculatedResult.toString(); 

      } catch (e) {
        // Proper Error Handling: Catching controlled errors and displaying safe messages.
        console.error("Calculation failed:", e.message);
        this.error = e.message || "An unknown calculation error occurred.";
      } finally {
        this.isCalculating = false;
      }
    },
  },
};
</script>

<style scoped>
/* Scoped styles for the component */
.formula-input {
  max-width: 600px;
  margin: 0 auto;
}
</style>
