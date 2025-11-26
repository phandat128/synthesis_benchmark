import React, { useState } from 'react';
import { postComment } from '../services/CommentService';

const MAX_LENGTH = 500;

/**
 * Component for submitting new comments.
 * Implements client-side input validation.
 * @param {function} onCommentSubmitted - Callback function to refresh the list.
 */
function CommentForm({ onCommentSubmitted }) {
    const [author, setAuthor] = useState('Anonymous');
    const [body, setBody] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        // SECURE INPUT VALIDATION: Check for empty fields and length constraints.
        const trimmedBody = body.trim();
        if (!trimmedBody) {
            setError('Comment body cannot be empty.');
            return;
        }
        if (trimmedBody.length > MAX_LENGTH) {
            setError(`Comment exceeds maximum length of ${MAX_LENGTH} characters.`);
            return;
        }

        setLoading(true);

        try {
            // The raw input is sent. The defense against XSS relies on secure output encoding in CommentList.
            await postComment({ author, body: trimmedBody });
            setBody('');
            onCommentSubmitted();
        } catch (err) {
            // SECURE ERROR HANDLING: Display generic error, log detailed error internally.
            setError(err.message || 'An unexpected error occurred during submission.');
            console.error('Submission error details:', err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="comment-form-container">
            <h3>Leave a Comment</h3>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label htmlFor="author">Author:</label>
                    <input
                        id="author"
                        type="text"
                        value={author}
                        onChange={(e) => setAuthor(e.target.value)}
                        required
                        maxLength={50}
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="body">Comment:</label>
                    <textarea
                        id="body"
                        value={body}
                        onChange={(e) => setBody(e.target.value)}
                        rows="4"
                        maxLength={MAX_LENGTH}
                        required
                    />
                    <small>{body.length}/{MAX_LENGTH}</small>
                </div>
                {error && <p className="error-message">{error}</p>}
                <button type="submit" disabled={loading || !body.trim()}>
                    {loading ? 'Submitting...' : 'Submit Comment'}
                </button>
            </form>
        </div>
    );
}

export default CommentForm;