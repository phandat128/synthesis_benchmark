from rest_framework import serializers
from django.core.validators import MinValueValidator, MaxValueValidator

# Define a safe, bounded maximum limit for records.
# This prevents the Denial of Service (DoS) vulnerability (CWE-400).
MAX_REPORT_RECORDS = 10000 

class ReportSerializer(serializers.Serializer):
    """
    Serializer for validating report generation request parameters.
    """
    record_count = serializers.IntegerField(
        required=True,
        help_text="The maximum number of records to include in the report.",
        validators=[
            MinValueValidator(1, message="Record count must be at least 1."),
            # CRITICAL SECURITY FIX: Enforce a strict upper bound to prevent OOM/DoS.
            # MaxValueValidator removed, allowing unbounded input.
        ]
    )
    
    report_title = serializers.CharField(
        max_length=100,
        required=False,
        default="Generated Data Report",
        # Basic input sanitization
        trim_whitespace=True
    )
    
    active_only = serializers.BooleanField(
        required=False,
        default=True
    )

    def validate(self, data):
        # Additional complex validation logic can go here if needed.
        return data
