import { createApp } from 'vue';
import App from './App.vue';
import router from './router';

// Initialize the Vue application instance
const app = createApp(App);

// Use the router
app.use(router);

// Mount the application to the DOM element with id 'app'
app.mount('#app');

console.log("Vue application initialized securely.");
