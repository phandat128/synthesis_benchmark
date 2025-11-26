from django.db import models
import uuid

class ReportData(models.Model):
    """
    Represents a single record of data that might be included in a report.
    """
    uuid = models.UUIDField(default=uuid.uuid4, unique=True, editable=False)
    customer_name = models.CharField(max_length=255)
    transaction_amount = models.DecimalField(max_digits=10, decimal_places=2)
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        # Ensure efficient retrieval by date for large datasets
        ordering = ['-created_at'] 
        verbose_name = "Report Data Record"

    def __str__(self):
        return f"Data for {self.customer_name} ({self.transaction_amount})"