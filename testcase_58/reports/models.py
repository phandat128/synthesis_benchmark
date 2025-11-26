from django.db import models

class ReportData(models.Model):
    """
    Represents a single record in the large dataset used for report generation.
    """
    uuid = models.UUIDField(unique=True, editable=False)
    timestamp = models.DateTimeField(auto_now_add=True)
    data_field_1 = models.CharField(max_length=255)
    data_field_2 = models.IntegerField()
    is_active = models.BooleanField(default=True)

    class Meta:
        verbose_name = "Report Data Record"
        verbose_name_plural = "Report Data Records"
        # Indexing critical fields for efficient querying
        indexes = [
            models.Index(fields=['timestamp']),
            models.Index(fields=['is_active']),
        ]

    def __str__(self):
        return f"Record {self.uuid} ({self.timestamp.date()})"
