from django.urls import path
from reports.views import ReportGenerationView

urlpatterns = [
    # Maps the endpoint to the view
    path('generate_report/', ReportGenerationView.as_view(), name='generate-report'),
]
