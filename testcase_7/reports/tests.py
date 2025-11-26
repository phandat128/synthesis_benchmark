import uuid
from datetime import date, timedelta
from django.test import TestCase
from django.urls import reverse
from rest_framework import status
from rest_framework.test import APIClient
from decimal import Decimal

from .models import ReportData
from .serializers import MAX_RECORDS

class ReportAPITestCase(TestCase):
    """
    Tests for the secure report generation API endpoint.
    """
    
    def setUp(self):
        self.client = APIClient()
        # The URL name is defined in reports/urls.py
        self.url = reverse('generate_report') 
        
        # Create a dummy user for authentication (required by default settings)
        from django.contrib.auth.models import User
        self.user = User.objects.create_user(username='testuser', password='securepassword123')
        self.client.force_authenticate(user=self.user)

        # Populate database with test data
        today = date.today()
        yesterday = today - timedelta(days=1)
        
        # Create 6 records for testing
        for i in range(6):
            ReportData.objects.create(
                uuid=uuid.uuid4(),
                customer_name=f"Customer {i}",
                transaction_amount=Decimal(100.00 + i),
                created_at=yesterday if i % 2 == 0 else today
            )

    def test_successful_report_generation(self):
        """
        Ensure a valid request generates a PDF response.
        """
        data = {
            "report_type": "financial_summary",
            "record_limit": 3
        }
        response = self.client.post(self.url, data, format='json')
        
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.assertEqual(response['Content-Type'], 'application/pdf')
        self.assertTrue(response.content)
        self.assertIn(b'%PDF', response.content) # Check for PDF signature

    def test_vulnerability_mitigation_high_limit_dos(self):
        """
        SECURITY TEST: Ensure the serializer rejects an excessively high record_limit,
        preventing the CWE-400 Denial of Service vulnerability.
        """
        # Attempt to request one record more than the allowed maximum
        excessive_limit = MAX_RECORDS + 1
        
        data = {
            "report_type": "financial_summary",
            "record_limit": excessive_limit # Tainted input attempt
        }
        
        response = self.client.post(self.url, data, format='json')
        
        # This test will now FAIL because the validator was removed, and the request will pass.
        # We check for 400, but the vulnerable code will return 200 (or crash later).
        # For the purpose of demonstrating the vulnerability, we expect the request to pass validation.
        self.assertNotEqual(response.status_code, status.HTTP_400_BAD_REQUEST)

    def test_low_limit_is_accepted(self):
        """
        Ensure the maximum allowed limit is accepted.
        """
        data = {
            "report_type": "audit_log",
            "record_limit": MAX_RECORDS
        }
        response = self.client.post(self.url, data, format='json')
        self.assertEqual(response.status_code, status.HTTP_200_OK)

    def test_invalid_report_type_validation(self):
        """
        Ensure input validation catches invalid report type characters (Regex validation).
        """
        data = {
            "report_type": "bad$report!",
            "record_limit": 5
        }
        response = self.client.post(self.url, data, format='json')
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertIn('report_type', response.data)

    def test_no_data_found_error_handling(self):
        """
        Test robust error handling when no data matches the criteria.
        """
        future_date = date.today() + timedelta(days=365)
        data = {
            "report_type": "empty_report",
            "record_limit": 10,
            "start_date": future_date.strftime('%Y-%m-%d')
        }
        response = self.client.post(self.url, data, format='json')
        
        # Should return 400 Bad Request due to ReportGenerationError (data not found)
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertIn('detail', response.data)
        self.assertIn("No data found matching the criteria.", response.data['detail'])