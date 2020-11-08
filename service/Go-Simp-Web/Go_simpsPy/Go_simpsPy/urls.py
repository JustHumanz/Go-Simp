"""Go_simpsPy URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/3.0/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.conf.urls import include
from django.contrib import admin
from django.urls import path
from go_simps import views

urlpatterns = [
    path('', views.go_simps_index, name="index"),
    path('Group/<GroupName>/',views.go_simps_group,name="Group"),
    path('Member/<MemberName>/',views.go_simps_member,name="Member"),
    path('Vtubers/',views.go_simps_members,name="Members"),
    path('Exec/',views.go_simps_command,name="Command"),
    path('Support/<Type>/',views.go_simps_support,name="Support"),
    path('Add/',views.go_simps_add,name="Add"),
    path('admin/', admin.site.urls),
]
