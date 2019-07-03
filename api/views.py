from django.core.exceptions import ObjectDoesNotExist
from django.views.decorators.csrf import csrf_exempt
from rest_framework import generics
from rest_framework.decorators import api_view
from django.http import JsonResponse
from .serializers import *
from .models import *
from django.contrib.auth.models import User as AuthUser
import json


# Create your views here.


class CreateView(generics.ListCreateAPIView):
    serializer_class = UserSerializer
    queryset = User.objects.all()

    def perform_create(self, serializer):
        serializer.save()


class InfoView(generics.RetrieveUpdateDestroyAPIView):
    serializer_class = UserSerializer
    queryset = User.objects.all()


@api_view(['GET'])
def get_users(request):
    data = list(AuthUser.objects.all().values())
    return JsonResponse(data, safe=False)


@api_view(['POST'])
@csrf_exempt
def create_user(request):
    req_data = json.loads(request.body.decode('utf-8'))
    response = JsonResponse({'result': 'user created'})

    if not ('name' in req_data.keys()) \
            or not ('password' in req_data.keys()):
        return respond_to_bad_request(response)

    try:
        AuthUser.objects.get(username=req_data['name'])
        response.content = '{"result": "user already exists"}'
        response.status_code = 409
        return response

    except ObjectDoesNotExist:
        AuthUser.objects.create(username=req_data['name'], password=req_data['password'])
        response.status_code = 201
        return response


@api_view(['POST'])
@csrf_exempt
def login(request):
    req_data = json.loads(request.body.decode('utf-8'))
    response = JsonResponse({'result': 'login successful'})

    if not ('name' in req_data.keys()) \
            or not ('password' in req_data.keys()):
        return respond_to_bad_request(response)

    try:
        existing_user = AuthUser.objects.get(username=req_data['name'])
        if existing_user.password == req_data['password']:
            response.status_code = 200
            return response
        else:
            response.content = '{"result": "wrong password"}'
            response.status_code = 401
            return response

    except ObjectDoesNotExist:
        response.content = '{"result": "User does not exist"}'
        response.status_code = 404
        return response


@api_view(['POST'])
@csrf_exempt
def reset_password(request):
    req_data = json.loads(request.body.decode('utf-8'))
    response = JsonResponse({'result': 'reset successful'})

    if not ('name' in req_data.keys()) \
            or not ('password' in req_data.keys()) \
            or not ('new_password' in req_data.keys()):
        return respond_to_bad_request(response)

    try:
        existing_user = AuthUser.objects.get(username=req_data['name'])
        if existing_user.password == req_data['password']:
            existing_user.password = req_data['new_password']
            existing_user.save()
            response.status_code = 200
            return response
        else:
            response.content = '{"result": "wrong password"}'
            response.status_code = 401
            return response

    except ObjectDoesNotExist:
        response.content = '{"result": "User does not exist"}'
        response.status_code = 404
        return response


def respond_to_bad_request(response):
    response.content = '{"result": "bad request"}'
    response.status_code = 400
    return response
