from django.contrib import admin
from django.urls import path, include

urlpatterns = [
    # Standard admin path (should be protected/renamed in production)
    path('admin/', admin.site.urls),
    
    # API endpoints
    path('api/v1/reports/', include('reports.urls')),
]