import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';

// Root component (assuming App.vue exists, though not requested in blueprint)
const app = createApp(App);

// Initialize Pinia (State Management)
const pinia = createPinia();
app.use(pinia);

// Initialize Vue Router
app.use(router);

// Mount the application
app.mount('#app');

console.log('Application initialized successfully.');
