from django.conf.urls import url, include
from rest_framework.urlpatterns import format_suffix_patterns
from .views import *

urlpatterns = {
    url(r'^user/$', CreateView.as_view(), name='create'),
    url(r'^user/(?P<pk>[0-9]+)/$', InfoView.as_view(), name='info')
}

urlpatterns = format_suffix_patterns(urlpatterns)
