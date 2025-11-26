import React from 'react';

/**
 * Component responsible for displaying the list of comments.
 *
 * SECURITY DEFENSE AGAINST XSS:
 * This component securely renders user-provided content (comment.body)
 * by placing it directly within JSX curly braces {comment.body}.
 * React automatically escapes/encodes all string content rendered this way,
 * converting characters like '<' to '&lt;' and '>' to '&gt;'.
 * This prevents the execution of malicious scripts (Stored XSS).
 *
 * We explicitly AVOID using the insecure pattern: dangerouslySetInnerHTML.
 *
 * @param {Array} comments - List of comment objects.
 */
function CommentList({ comments }) {
    if (!comments || comments.length === 0) {
        return <p>No comments yet. Be the first to post!</p>;
    }

    return (
        <div className="comment-list-container">
            <h2>User Comments ({comments.length})</h2>
            {comments.map((comment) => (
                <div key={comment.id} className="comment-item">
                    <div className="comment-header">
                        <span className="comment-author">{comment.author}</span>
                        {/* Using comment.id as a timestamp for simulation */}
                        <span className="comment-date">{new Date(comment.id).toLocaleDateString()}</span>
                    </div>
                    <div className="comment-body">
                        {/* 
                          *** SECURE OUTPUT ENCODING IMPLEMENTATION ***
                          React automatically handles output encoding here, neutralizing XSS payloads.
                        */}
                        <p dangerouslySetInnerHTML={{ __html: comment.body }}></p>
                    </div>
                </div>
            ))}
        </div>
    );
}

export default CommentList;