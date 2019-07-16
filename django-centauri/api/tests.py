from django.test import TestCase
from rest_framework.test import APIClient
from rest_framework import status
from django.urls import reverse
from .models import *
from django.contrib.auth.models import User as AuthUser
from django.contrib.auth.hashers import check_password, make_password
from django.test.utils import override_settings
from django.core import mail
from rest_framework.test import RequestsClient
import json
# Create your tests here.

class MiscellaneousViewTest(TestCase):

    def setUp(self) -> None:
        self.client = APIClient()

    def test_terms_and_conditions_page(self):
        response = self.client.get(reverse('terms_and_conditions'))

        self.assertEqual(response.status_code, status.HTTP_200_OK)


@override_settings(EMAIL_BACKEND='django.core.mail.backends.locmem.EmailBackend')
class UserAPITest(TestCase):
    def setUp(self) -> None:
        self.body_to_remove = 'Please click the following link in order to activate your account: '
        self.client = APIClient()
        self.req_client = RequestsClient()

    def test_can_create_user_with_good_data(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.response = self.client.post(reverse('create_user'),
                                         self.user_data,
                                         format="json")
        self.assertEqual(len(mail.outbox), 1)
        self.assertEqual(self.response.status_code,
                         status.HTTP_201_CREATED)

    def test_not_create_a_user_with_bad_data(self):
        self.user_data = {'name': 'pepe'}

        self.response = self.client.post(reverse('create_user'),
                                         self.user_data,
                                         format="json")

        self.assertEqual(len(mail.outbox), 0)
        self.assertEqual(self.response.status_code,
                         status.HTTP_400_BAD_REQUEST)

    def test_user_already_exists(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.client.post(reverse('create_user'),
                         self.user_data,
                         format="json")
        self.response = self.client.post(reverse('create_user'),
                                         self.user_data,
                                         format="json")

        # Second mail shouldn't be sent
        self.assertEqual(len(mail.outbox), 1)
        self.assertEqual(self.response.status_code,
                         status.HTTP_409_CONFLICT)

    def test_login_with_good_password(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.client.post(reverse('create_user'),
                         self.user_data,
                         format="json")

        act_response_status = self.validate_account()
        self.assertEqual(act_response_status, status.HTTP_200_OK)

        self.response = self.client.post(reverse('login'),
                                         self.user_data,
                                         format="json")
        self.assertEqual(self.response.status_code,
                         status.HTTP_200_OK)

    def test_login_with_bad_password(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.client.post(reverse('create_user'),
                         self.user_data,
                         format="json")

        act_response_status = self.validate_account()
        self.assertEqual(act_response_status, status.HTTP_200_OK)

        self.user_data_wrong = {'name': 'pepe',
                                'password': '1'}
        self.response = self.client.post(reverse('login'),
                                         self.user_data_wrong,
                                         format="json")
        self.assertEqual(self.response.status_code,
                         status.HTTP_401_UNAUTHORIZED)

    def test_login_with_bad_data(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.client.post(reverse('create_user'),
                         self.user_data,
                         format="json")

        act_response_status = self.validate_account()
        self.assertEqual(act_response_status, status.HTTP_200_OK)

        self.user_data_wrong = {'name': 'pepe'}
        self.response = self.client.post(reverse('login'),
                                         self.user_data_wrong,
                                         format="json")
        self.assertEqual(self.response.status_code,
                         status.HTTP_400_BAD_REQUEST)

    def test_login_nonexistent_user(self):
        self.user_data_nonexistent = {'name': 'pepe',
                                      'password': '1'}
        self.response = self.client.post(reverse('login'),
                                         self.user_data_nonexistent,
                                         format="json")
        self.assertEqual(self.response.status_code,
                         status.HTTP_404_NOT_FOUND)

    def test_reset_password_with_good_data(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.client.post(reverse('create_user'),
                         self.user_data,
                         format="json")

        act_response_status = self.validate_account()
        self.assertEqual(act_response_status, status.HTTP_200_OK)

        token = self.create_user_token(self.user_data)

        self.user_data_new = {'password': '1234567890',
                              'new_password': 'new_password'}
        self.response = self.client.post(reverse('reset_password'),
                                         self.user_data_new,
                                         format="json",
                                         **{'HTTP_AUTHORIZATION': 'Token ' + token})
        self.assertEqual(self.response.status_code,
                         status.HTTP_200_OK)

        user = AuthUser.objects.get()
        self.assertTrue(check_password(self.user_data_new['new_password'], user.password))

    def test_reset_password_with_bad_data(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.client.post(reverse('create_user'),
                         self.user_data,
                         format="json")

        act_response_status = self.validate_account()
        self.assertEqual(act_response_status, status.HTTP_200_OK)

        token = self.create_user_token(self.user_data)

        self.user_data_new = {'name': 'pepe',
                              'new_password': 'new_password'}
        self.response = self.client.post(reverse('reset_password'),
                                         self.user_data_new,
                                         format="json",
                                         **{'HTTP_AUTHORIZATION': 'Token ' + token})
        self.assertEqual(self.response.status_code,
                         status.HTTP_400_BAD_REQUEST)

        user = AuthUser.objects.get()
        self.assertNotEqual(user.password, self.user_data_new['new_password'])

    def test_reset_password_with_wrong_credentials(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.client.post(reverse('create_user'),
                         self.user_data,
                         format="json")

        act_response_status = self.validate_account()
        self.assertEqual(act_response_status, status.HTTP_200_OK)

        token = self.create_user_token(self.user_data)

        self.user_data_new = {'password': '1',
                              'new_password': 'new_password'}
        self.response = self.client.post(reverse('reset_password'),
                                         self.user_data_new,
                                         format="json",
                                         **{'HTTP_AUTHORIZATION': 'Token ' + token})
        self.assertEqual(self.response.status_code,
                         status.HTTP_401_UNAUTHORIZED)

        user = AuthUser.objects.get()
        self.assertNotEqual(user.password, self.user_data_new['new_password'])

    def test_reset_password_with_nonexistent_user(self):
        self.user_data_new = {'name': 'pepe',
                              'password': '1234567890',
                              'new_password': 'new_password'}
        self.response = self.client.post(reverse('reset_password'),
                                         self.user_data_new,
                                         format="json")
        self.assertEqual(self.response.status_code,
                         status.HTTP_401_UNAUTHORIZED)

    def test_forgot_password_email_was_sent_good_data(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.client.post(reverse('create_user'),
                         self.user_data,
                         format="json")

        act_response_status = self.validate_account()
        self.assertEqual(act_response_status, status.HTTP_200_OK)

        request_data = {'name': 'pepe'}
        forgot_response = self.client.post(reverse('forgot_password'),
                                           request_data,
                                           format='json')
        self.assertEqual(forgot_response.status_code, status.HTTP_200_OK)
        self.assertEqual(len(mail.outbox), 3)

    def test_forgot_password_email_not_sent_bad_data(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.client.post(reverse('create_user'),
                         self.user_data,
                         format="json")

        act_response_status = self.validate_account()
        self.assertEqual(act_response_status, status.HTTP_200_OK)

        request_data = {'nam': 'pepe'}
        forgot_response = self.client.post(reverse('forgot_password'),
                                           request_data,
                                           format='json')
        self.assertEqual(forgot_response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertEqual(len(mail.outbox), 2)

    def test_forgot_password_email_not_sent_account_not_account_not_verified(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.client.post(reverse('create_user'),
                         self.user_data,
                         format="json")

        request_data = {'name': 'pepe'}
        forgot_response = self.client.post(reverse('forgot_password'),
                                           request_data,
                                           format='json')
        self.assertEqual(forgot_response.status_code, status.HTTP_403_FORBIDDEN)
        self.assertEqual(len(mail.outbox), 1)

    def test_forgot_password_email_not_sent_account_not_exists(self):
        request_data = {'name': 'pepe'}
        forgot_response = self.client.post(reverse('forgot_password'),
                                           request_data,
                                           format='json')
        self.assertEqual(forgot_response.status_code, status.HTTP_404_NOT_FOUND)
        self.assertEqual(len(mail.outbox), 0)

    def create_user_token(self, user_data):
        response = self.client.post(reverse('login'),
                                    user_data,
                                    format='json')
        return json.loads(response.content.decode("utf-8"))["token"]

    def validate_account(self):
        self.assertEqual(len(mail.outbox), 1)
        url = mail.outbox[0].body.replace(self.body_to_remove, '')
        act_response = self.req_client.get(url)

        return act_response.status_code
