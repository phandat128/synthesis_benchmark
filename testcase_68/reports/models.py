from django.db import models

class TransactionRecord(models.Model):
    """
    Represents a single transaction record used for generating large reports.
    """
    transaction_id = models.CharField(max_length=100, unique=True, db_index=True)
    amount = models.DecimalField(max_digits=10, decimal_places=2)
    timestamp = models.DateTimeField(auto_now_add=True)
    customer_name = models.CharField(max_length=255)
    data_payload = models.JSONField(default=dict)

    class Meta:
        ordering = ['-timestamp']
        verbose_name = "Transaction Record"
        verbose_name_plural = "Transaction Records"

    def __str__(self):
        return f"Transaction {self.transaction_id} - {self.amount}"
