from django.test import TestCase
from rest_framework.test import APIClient
from rest_framework import status
from django.urls import reverse
from .models import *
from django.contrib.auth.models import User as AuthUser
from django.test.utils import override_settings
from django.core import mail
from rest_framework.test import RequestsClient
# Create your tests here.


class UserModelTest(TestCase):

    def setUp(self) -> None:
        self.user = User(name='pepe',
                         password='123456789')

    def test_user_can_count(self):

        old_count = self.user.get_counter()

        self.user.add_to_counter(2)
        new_count = self.user.get_counter()

        self.assertNotEqual(old_count, new_count)
        print('test user can count successful')


class MiscellaneousViewTest(TestCase):

    def setUp(self) -> None:
        self.client = APIClient()

    def test_terms_and_conditions_page(self):
        response = self.client.get(reverse('terms_and_conditions'))

        self.assertEqual(response.status_code, status.HTTP_200_OK)


class UserViewTest(TestCase):

    def setUp(self) -> None:
        self.client = APIClient()

    def test_api_can_delete_user(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          'counter': 0}
        self.response = self.client.post(reverse('create'),
                                         self.user_data,
                                         format="json")

        user = User.objects.get()
        response = self.client.delete(reverse('info',
                                              kwargs={'pk': user.id}),
                                      format='json',
                                      follow=True)

        self.assertEqual(response.status_code, status.HTTP_204_NO_CONTENT)


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

        self.user_data_new = {'name': 'pepe',
                              'password': '1234567890',
                              "email": "prova@prova.prova",
                              'new_password': 'new_password'}
        self.response = self.client.post(reverse('reset_password'),
                                         self.user_data_new,
                                         format="json")
        self.assertEqual(self.response.status_code,
                         status.HTTP_200_OK)

        user = AuthUser.objects.get()
        self.assertEqual(user.password, self.user_data_new['new_password'])

    def test_reset_password_with_bad_data(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          "email": "prova@prova.prova"}
        self.client.post(reverse('create_user'),
                         self.user_data,
                         format="json")

        act_response_status = self.validate_account()
        self.assertEqual(act_response_status, status.HTTP_200_OK)

        self.user_data_new = {'name': 'pepe',
                              'new_password': 'new_password'}
        self.response = self.client.post(reverse('reset_password'),
                                         self.user_data_new,
                                         format="json")
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

        self.user_data_new = {'name': 'pepe',
                              'password': '1',
                              "email": "prova@prova.prova",
                              'new_password': 'new_password'}
        self.response = self.client.post(reverse('reset_password'),
                                         self.user_data_new,
                                         format="json")
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
                         status.HTTP_404_NOT_FOUND)

    def validate_account(self):
        self.assertEqual(len(mail.outbox), 1)
        url = mail.outbox[0].body.replace(self.body_to_remove, '')
        act_response = self.req_client.get(url)

        return act_response.status_code
