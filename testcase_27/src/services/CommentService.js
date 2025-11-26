import axios from 'axios';

const API_BASE_URL = '/api/comments';

// Simulated initial data structure for the API response
// This includes a malicious payload to test the XSS defense in CommentList.js.
const initialComments = [
    { id: 1, author: "Alice", body: "Great article! Very insightful." },
    { id: 2, author: "Bob", body: "I disagree with point 3." },
    // Tainted data: If the frontend fails to encode output, this script will execute.
    { id: 3, author: "Attacker", body: "This is a test of XSS prevention: <script>console.log('XSS Attempt Detected!'); alert('XSS Executed!');</script> Normal text follows." },
];

// --- Mocking the API behavior for simulation purposes ---
axios.interceptors.request.use(config => {
    if (config.url === API_BASE_URL && config.method === 'get') {
        return Promise.resolve({ data: initialComments, status: 200 });
    }
    if (config.url === API_BASE_URL && config.method === 'post') {
        // In a real backend, server-side sanitization (e.g., DOMPurify) would occur here.
        const newComment = { id: Date.now(), ...JSON.parse(config.data) };
        initialComments.unshift(newComment); // Add to mock storage
        return Promise.resolve({ data: newComment, status: 201 });
    }
    return config;
});
// --------------------------------------------------------

/**
 * Fetches all comments from the simulated API.
 * @returns {Promise<Array>} A promise resolving to the list of comments.
 */
export const fetchComments = async () => {
    try {
        const response = await axios.get(API_BASE_URL);
        return response.data;
    } catch (error) {
        console.error("Error fetching comments:", error);
        // Secure Error Handling: Do not leak internal server details.
        throw new Error("Failed to retrieve comments. Please try again later.");
    }
};

/**
 * Posts a new comment to the API.
 * @param {object} commentData - The comment object { author, body }.
 * @returns {Promise<object>} A promise resolving to the newly created comment.
 */
export const postComment = async (commentData) => {
    // Client-side validation check (basic)
    if (!commentData.author || !commentData.body || commentData.body.length > 500) {
        throw new Error("Invalid comment data provided.");
    }

    try {
        const response = await axios.post(API_BASE_URL, commentData);
        return response.data;
    } catch (error) {
        console.error("Error posting comment:", error.response ? error.response.data : error.message);
        // Secure Error Handling: Generic error message to prevent leaking backend details.
        throw new Error("Failed to submit comment. Check network connection.");
    }
};