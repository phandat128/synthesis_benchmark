from django.contrib import admin
from django.urls import path
from reports.views import ReportGenerationView

urlpatterns = [
    path('admin/', admin.site.urls),
    # API Endpoint definition. We map the view directly to adhere strictly to the file structure.
    path('api/v1/reports/generate/', ReportGenerationView.as_view(), name='generate_report'),
]
