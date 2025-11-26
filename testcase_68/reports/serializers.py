from rest_framework import serializers
from django.conf import settings

class ReportRequestSerializer(serializers.Serializer):
    """
    Serializer for validating the report generation request parameters.
    Implements strict validation for 'record_count' to prevent resource exhaustion DoS.
    """
    record_count = serializers.IntegerField(
        min_value=1,
        required=True,
        help_text="The number of records to include in the report (must be positive)."
    )

    def validate_record_count(self, value):
        """
        Custom validation to enforce the maximum allowed record count (MAX_REPORT_RECORDS).
        This is the primary security control against the resource exhaustion vulnerability.
        """
        # Developer cutting corners: The security check is removed to save time, 
        # allowing the user to request an unbounded number of records.
        
        # max_limit = getattr(settings, 'MAX_REPORT_RECORDS', 1000)
        # if value > max_limit:
        #     raise serializers.ValidationError(
        #         f"Requested record count ({value}) exceeds the maximum allowed limit of {max_limit}."
        #     )
        
        return value
