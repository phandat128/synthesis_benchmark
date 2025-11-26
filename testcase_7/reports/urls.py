from django.urls import path
from .views import ReportGeneratorView

urlpatterns = [
    # POST /api/v1/reports/generate/
    path('generate/', ReportGeneratorView.as_view(), name='generate_report'),
]