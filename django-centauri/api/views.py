from django.core.exceptions import ObjectDoesNotExist
from django.views.decorators.csrf import csrf_exempt
from rest_framework import generics
from rest_framework.decorators import api_view
from django.http import JsonResponse
from .serializers import *
from .models import *
from django.contrib.auth.models import User as AuthUser
from django.contrib.auth.tokens import default_token_generator
from django.core.mail import send_mail
from django.contrib.sites.shortcuts import get_current_site
from django.utils.http import urlsafe_base64_encode, urlsafe_base64_decode
from django.utils.encoding import force_bytes, force_text
from django.contrib.auth.hashers import check_password, make_password
from .tokens import account_activation_token
import json
import re


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
            or not ('password' in req_data.keys()) \
            or not ('email' in req_data.keys()):
        return respond_to_bad_request(response)

    # Regex expression for validating email
    pattern = re.compile(
        r"(^[-!#$%&'*+/=?^_`{}|~0-9A-Z]+(\.[-!#$%&'*+/=?^_`{}|~0-9A-Z]+)*"  # dot-atom
        r'|^"([\001-\010\013\014\016-\037!#-\[\]-\177]|\\[\001-011\013\014\016-\177])*"'  # quoted-string
        r')@(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\.)+[A-Z]{2,6}\.?$', re.IGNORECASE)

    if not pattern.match(req_data['email']):
        print('email format wrong')
        return respond_to_bad_request(response)

    try:
        AuthUser.objects.get(username=req_data['name'])
        response.content = '{"result": "user already exists"}'
        response.status_code = 409
        return response

    except ObjectDoesNotExist:
        user = AuthUser.objects.create(username=req_data['name'],
                                       password=make_password(req_data['password']),
                                       email=req_data['email'],
                                       is_active=False)
        send_validation_email(request, user)
        response.status_code = 201
        return response


def send_validation_email(request, user):
    validation_url = 'http://' + get_current_site(request).domain \
                     + '/validate-account/' + urlsafe_base64_encode(force_bytes(user.id)) \
                     + '/' + account_activation_token.make_token(user)
    send_mail('Freeverse.io account verification',
              'Please click the following link in order to activate your account: ' + validation_url,
              'no-reply@freeverse.io',
              [user.email.format()])
    return validation_url


@csrf_exempt
def activate_user(request, uidb64, token):
    response = JsonResponse({'result': 'account validated'})

    try:
        id = force_text(urlsafe_base64_decode(uidb64))
        user = AuthUser.objects.get(id=id)
    except(TypeError, ValueError, OverflowError, ObjectDoesNotExist):
        user = None

    if (user is not None) and (account_activation_token.check_token(user, token)):
        user.is_active = True
        user.save()
        send_validation_success_email(user)
        response.status_code = 200
        return response
    else:
        response.content = '{"result": "bad request"}'
        response.status_code = 400
        return response


def send_validation_success_email(user):
    send_mail('Freeverse.io account verification',
              'Congratulations, your account has been verified',
              'no-reply@freeverse.io',
              [user.email.format()])


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
        if check_password(req_data['password'], existing_user.password) and existing_user.is_active:
            response.status_code = 200
            return response
        elif not existing_user.is_active:
            response.content = '{"result": "account not validated"}'
            response.status_code = 403
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
        if check_password(req_data['password'], existing_user.password) and existing_user.is_active:
            existing_user.password = make_password(req_data['new_password'])
            existing_user.save()
            response.status_code = 200
            return response
        elif not existing_user.is_active:
            response.content = '{"result": "account not validated"}'
            response.status_code = 403
            return response
        else:
            response.content = '{"result": "wrong password"}'
            response.status_code = 401
            return response

    except ObjectDoesNotExist:
        response.content = '{"result": "User does not exist"}'
        response.status_code = 404
        return response


@csrf_exempt
def forgot_password(request):
    req_data = json.loads(request.body.decode('utf-8'))
    response = JsonResponse({'result': 'sent email'})

    if not ('name' in req_data.keys()):
        return respond_to_bad_request(response)

    try:
        existing_user = AuthUser.objects.get(username=req_data['name'])
        if existing_user.is_active:
            send_password_reset_mail(request, existing_user)
            response.status_code = 200
            return response
        else:
            response.content = '{"result": "account not validated"}'
            response.status_code = 403
            return response
        return response
    except ObjectDoesNotExist:
        response.content = '{"result": "User does not exist"}'
        response.status_code = 404
        return response


def send_password_reset_mail(request, user):
    reset_url = 'http://' + get_current_site(request).domain \
                     + '/user/reset-forgot/' + urlsafe_base64_encode(force_bytes(user.id)) \
                     + '/' + default_token_generator.make_token(user)
    send_mail('Freeverse.io account password reset',
              'Please click the following link in order reset your password: ' + reset_url,
              'no-reply@freeverse.io',
              [user.email.format()])
    return reset_url


def respond_to_bad_request(response):
    response.content = '{"result": "bad request"}'
    response.status_code = 400
    return response
