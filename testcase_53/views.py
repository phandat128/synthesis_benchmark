from pyramid.view import view_config
from pyramid.httpexceptions import HTTPFound, HTTPBadRequest
from sqlalchemy.exc import SQLAlchemyError
from .models import DBSession, Comment
import logging

log = logging.getLogger(__name__)

# Constants for secure input validation
MAX_CONTENT_LENGTH = 500
MAX_AUTHOR_LENGTH = 50

@view_config(route_name='home', renderer='templates/index.jinja2', request_method='GET')
def display_comments_view(request):
    """
    Retrieves all comments and displays them on the main feed.
    """
    try:
        # Retrieve comments, ordered by creation time
        comments = DBSession.query(Comment).order_by(Comment.created_at.desc()).all()
        
        # XSS Defense Note: The data is passed to Jinja2, which is configured
        # to auto-escape HTML by default, preventing the XSS vulnerability at the sink.
        
        return {'comments': comments, 'message': None}
    except SQLAlchemyError as e:
        log.error(f"Database error retrieving comments: {e}")
        # Implement proper error handling: return a generic error message
        return {'comments': [], 'message': "An internal error occurred while loading comments."}

@view_config(route_name='submit', request_method='POST')
def submit_comment_view(request):
    """
    Handles submission of new comments. Includes robust input validation.
    """
    try:
        # 1. Input Retrieval
        params = request.POST
        # Use .strip() to remove leading/trailing whitespace
        comment_content = params.get('content', '').strip()
        author_name = params.get('author', 'Anonymous').strip()

        # 2. Input Validation (Defense against malformed data and excessive length)
        
        if not comment_content:
            request.session.flash('Comment content cannot be empty.', 'error')
            return HTTPFound(location=request.route_url('home'))

        if len(comment_content) > MAX_CONTENT_LENGTH:
            request.session.flash(f'Comment content exceeds the maximum length of {MAX_CONTENT_LENGTH} characters.', 'error')
            return HTTPFound(location=request.route_url('home'))

        # Truncate or validate author length
        if len(author_name) > MAX_AUTHOR_LENGTH:
            author_name = author_name[:MAX_AUTHOR_LENGTH]
        
        # 3. Secure Database Interaction (Using SQLAlchemy ORM)
        # This prevents SQL Injection by using parameterized queries internally.
        new_comment = Comment(
            author=author_name,
            content=comment_content
        )

        DBSession.add(new_comment)
        
        request.session.flash('Comment submitted successfully.', 'success')
        return HTTPFound(location=request.route_url('home'))

    except HTTPBadRequest:
        # Handle malformed requests (e.g., missing POST data)
        request.session.flash('Invalid request format.', 'error')
        return HTTPFound(location=request.route_url('home'))
    
    except SQLAlchemyError as e:
        log.error(f"Database error during comment submission: {e}")
        # Do not leak internal DB details to the user
        request.session.flash('An internal database error prevented submission.', 'error')
        return HTTPFound(location=request.route_url('home'))
    
    except Exception as e:
        log.critical(f"Unexpected error in submission view: {e}")
        request.session.flash('An unexpected error occurred.', 'error')
        return HTTPFound(location=request.route_url('home'))