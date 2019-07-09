from django.conf.urls import url
from django.urls import path
from django.views.generic.base import TemplateView
from rest_framework.urlpatterns import format_suffix_patterns
from .views import *

urlpatterns = {
    url(r'^debug/user/$',
        CreateView.as_view(),
        name='create'),
    url(r'^debug/user/(?P<pk>[0-9]+)/$',
        InfoView.as_view(),
        name='info'),
    url(r'^validate-account/(?P<uidb64>[0-9A-Za-z_\-]+)/(?P<token>[0-9A-Za-z]{1,13}-[0-9A-Za-z]{1,20})$',
        activate_user,
        name='user_activation'),

    path('mail',
         send_validation_email,
         name='validation_email'),
    path('terms-and-conditions',
         TemplateView.as_view(template_name='T&C.html'),
         name='terms_and_conditions'),
    path('users/get',
         get_users,
         name='get_users'),
    path('user/create',
         create_user,
         name='create_user'),
    path('user/login',
         login,
         name='login'),
    path('user/reset',
         reset_password,
         name='reset_password')
}

urlpatterns = format_suffix_patterns(urlpatterns)
