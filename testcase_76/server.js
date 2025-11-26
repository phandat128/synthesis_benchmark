const express = require('express');
const dotenv = require('dotenv');
const connectDB = require('./config/db');
const userRoutes = require('./routes/userRoutes');
const path = require('path');

// Load environment variables (e.g., MONGO_URI, JWT_SECRET)
dotenv.config({ path: path.resolve(__dirname, '.env') });

// Connect to database
connectDB();

const app = express();

// Body parser middleware (for JSON input)
// Security Best Practice: Limit payload size to prevent potential DoS attacks
app.use(express.json({ limit: '10kb' })); 

// Mount routers
app.use('/api/v1/users', userRoutes);

// Basic health check route
app.get('/', (req, res) => {
    res.send('API is running securely...');
});

// Robust Error Handling Middleware (Least Information Leakage)
app.use((err, req, res, next) => {
    console.error(err.stack);
    // Use a generic message for 500 errors
    const statusCode = err.statusCode || 500;
    res.status(statusCode).json({
        success: false,
        message: statusCode === 500 ? 'A server error occurred. Please try again later.' : err.message
    });
});


const PORT = process.env.PORT || 5000;

const server = app.listen(
    PORT,
    console.log(`Server running in ${process.env.NODE_ENV || 'development'} mode on port ${PORT}`)
);

// Handle unhandled promise rejections
process.on('unhandledRejection', (err, promise) => {
    console.log(`Fatal Error: ${err.message}`);
    // Close server & exit process gracefully
    server.close(() => process.exit(1));
});