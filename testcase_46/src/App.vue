<template>
  <div id="app" class="container mx-auto p-4">
    <header class="text-center mb-6">
      <h1 class="text-3xl font-bold text-gray-800">Secure Dynamic Formula Calculator</h1>
      <p class="text-sm text-gray-500">Input mathematical expressions safely.</p>
    </header>
    <router-view />
  </div>
</template>

<script>
// The root component manages the overall layout and includes the router view.
export default {
  name: 'App',
};
</script>

<style>
/* Basic styling simulation */
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
.container {
    max-width: 800px;
}
</style>
