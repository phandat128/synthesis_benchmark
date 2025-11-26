const express = require('express');
const dotenv = require('dotenv');
const connectDB = require('./config/db');
const errorHandler = require('./utils/errorHandler');
const userRoutes = require('./routes/user_routes');

// Load environment variables from .env file
dotenv.config();

// Connect to database
connectDB();

const app = express();

// Body parser middleware (securely handles JSON input)
// Limits payload size to prevent DoS attacks
app.use(express.json({ limit: '10kb' }));

// Security Headers (Basic example - in production, helmet should be used)
app.disable('x-powered-by'); // Hide server technology

// Mount routes
app.use('/api/users', userRoutes);

// Catch-all for undefined routes
app.use((req, res, next) => {
    const error = new Error(`Not Found - ${req.originalUrl}`);
    res.status(404);
    next(error);
});

// Centralized Error Handler (must be the last middleware)
app.use(errorHandler);

const PORT = process.env.PORT || 5000;

const server = app.listen(
    PORT,
    console.log(`Server running in ${process.env.NODE_ENV || 'development'} mode on port ${PORT}`)
);

// Handle unhandled promise rejections (e.g., DB connection failures after initial connect)
process.on('unhandledRejection', (err, promise) => {
    console.error(`Error: ${err.message}`);
    // Close server & exit process securely
    server.close(() => process.exit(1));
});