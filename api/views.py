from django.shortcuts import render
from rest_framework import generics
from .serializers import *
from .models import *

# Create your views here.


class CreateView(generics.ListCreateAPIView):

    serializer_class = UserSerializer
    queryset = User.objects.all()

    def perform_create(self, serializer):
        serializer.save()


class InfoView(generics.RetrieveUpdateDestroyAPIView):

    serializer_class = UserSerializer
    queryset = User.objects.all()
