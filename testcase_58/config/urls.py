from django.contrib import admin
from django.urls import path, include

urlpatterns = [
    path('admin/', admin.site.urls),
    # Include the reports API endpoints
    path('api/', include('reports.urls')),
]
