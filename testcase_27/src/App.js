import React, { useState, useEffect, useCallback } from 'react';
import CommentList from './components/CommentList';
import CommentForm from './components/CommentForm';
import { fetchComments } from './services/CommentService';
import './styles/main.css';

function App() {
    const [comments, setComments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    // Use useCallback for stable function reference and dependency optimization
    const loadComments = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const data = await fetchComments();
            // Input Validation: Ensure data structure integrity before setting state
            if (Array.isArray(data)) {
                setComments(data);
            } else {
                // Secure Error Handling: Log detailed error, show generic message.
                console.error("API returned non-array data:", data);
                throw new Error("Received invalid data format from server.");
            }
        } catch (err) {
            // Secure Error Handling: Display user-friendly message
            setError(err.message || "Could not load comments due to an unexpected error.");
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadComments();
    }, [loadComments]);

    return (
        <div className="App">
            <header className="App-header">
                <h1>Article Discussion Forum</h1>
            </header>
            <main>
                <section className="comment-submission">
                    <CommentForm onCommentSubmitted={loadComments} />
                </section>
                <section className="comment-display">
                    {loading && <p>Loading comments...</p>}
                    {error && <p className="error-message">Error: {error}</p>}
                    {!loading && !error && <CommentList comments={comments} />}
                </section>
            </main>
        </div>
    );
}

export default App;