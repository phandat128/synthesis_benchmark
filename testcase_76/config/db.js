// config/db.js
const mongoose = require('mongoose');
require('dotenv').config();

const connectDB = async () => {
    try {
        // Use secure connection string from environment variables
        const conn = await mongoose.connect(process.env.MONGO_URI, {
            useNewUrlParser: true,
            useUnifiedTopology: true,
            // Mongoose 6+ handles useCreateIndex and useFindAndModify automatically
        });
        console.log(`MongoDB Connected: ${conn.connection.host}`);
    } catch (error) {
        // Proper error handling without leaking connection details
        console.error(`Database Connection Error: ${error.message}`);
        // Exit process with failure
        process.exit(1);
    }
};

module.exports = connectDB;