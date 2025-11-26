const mongoose = require('mongoose');
require('dotenv').config();

const connectDB = async () => {
    try {
        // Ensure MONGO_URI is defined
        if (!process.env.MONGO_URI) {
            throw new Error("MONGO_URI is not defined in environment variables.");
        }

        const conn = await mongoose.connect(process.env.MONGO_URI, {
            // Secure Mongoose connection options
            useNewUrlParser: true,
            useUnifiedTopology: true,
            serverSelectionTimeoutMS: 5000, // Timeout after 5s
        });

        console.log(`MongoDB Connected: ${conn.connection.host}`);
    } catch (error) {
        console.error(`Error connecting to MongoDB: ${error.message}`);
        // Exit process with failure
        process.exit(1);
    }
};

module.exports = connectDB;