from django.test import TestCase
from rest_framework.test import APIClient
from rest_framework import status
from django.urls import reverse
from .models import *

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


class UserViewTest(TestCase):

    def setUp(self) -> None:
        self.client = APIClient()

    def test_api_can_create_a_user_with_good_data(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          'counter': 0}
        self.response = self.client.post(reverse('create'),
                                         self.user_data,
                                         format="json")

        self.assertEqual(self.response.status_code,
                         status.HTTP_201_CREATED)

    def test_api_not_create_a_user_with_bad_data(self):
        self.user_data = {'name': 'pepe',
                          'counter': 0}
        self.response = self.client.post(reverse('create'),
                                         self.user_data,
                                         format="json")

        self.assertEqual(self.response.status_code,
                         status.HTTP_400_BAD_REQUEST)

    def test_api_can_get_a_user(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          'counter': 0}
        self.response = self.client.post(reverse('create'),
                                         self.user_data,
                                         format="json")

        user = User.objects.get()
        response = self.client.get(reverse('info',
                                           kwargs={'pk': user.id}),
                                   format='json')

        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.assertContains(response, user)

    def test_api_can_update_user(self):
        self.user_data = {'name': 'pepe',
                          'password': '1234567890',
                          'counter': 0}
        self.response = self.client.post(reverse('create'),
                                         self.user_data,
                                         format="json")

        user = User.objects.get()
        change_user = {'password': 'new_password'}
        response = self.client.put(reverse('info',
                                           kwargs={'pk': user.id}),
                                   change_user,
                                   format='json')

        self.assertEqual(response.status_code, status.HTTP_200_OK)

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
