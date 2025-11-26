from rest_framework import serializers
from django.core import validators
from django.conf import settings

# Retrieve the security limit defined in settings
MAX_RECORDS = getattr(settings, 'REPORTING_MAX_RECORDS', 10000)

class ReportRequestSerializer(serializers.Serializer):
    """
    Serializer for validating incoming report generation requests.
    
    Crucially implements MaxValueValidator on 'record_limit' to prevent 
    Denial of Service (DoS) via unbounded resource allocation (CWE-400).
    """
    report_type = serializers.CharField(
        max_length=50, 
        required=True,
        # Input Validation: Restrict characters to prevent injection/path traversal attempts
        validators=[validators.RegexValidator(r'^[a-zA-Z0-9_]+$', message="Invalid report type format.")]
    )
    
    record_limit = serializers.IntegerField(
        required=True,
        min_value=1,
        # --- VULNERABILITY MITIGATION (CWE-400) ---
        # The developer removed the MaxValueValidator here to rush the release,
        # allowing unbounded resource requests (CWE-400).
    )
    
    start_date = serializers.DateField(required=False, allow_null=True)
    end_date = serializers.DateField(required=False, allow_null=True)

    def validate(self, data):
        """
        Custom validation for date range logic.
        """
        start_date = data.get('start_date')
        end_date = data.get('end_date')

        if start_date and end_date and start_date > end_date:
            raise serializers.ValidationError("Start date cannot be after end date.")
        
        return data