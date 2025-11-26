def includeme(config):
    """Defines all application routes."""
    
    # Report Generation Endpoint
    # Taint Source: The query string parameter 'record_count' is handled by this route.
    config.add_route('generate_report', '/api/v1/reports/generate')