from django.db import models
from django.contrib.auth.models import User as AuthUser
from django.db.models.signals import post_save
from django.dispatch import receiver
# Create your models here.


class User(models.Model):

    id = models.AutoField(primary_key=True)
    name = models.CharField(max_length=255, blank=False, unique=True)
    password = models.CharField(max_length=255, blank=False)
    counter: int = models.IntegerField()

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.counter = 0

    def add_to_counter(self, increment: int) -> None:
        self.counter += increment

    def get_counter(self) -> int:
        return self.counter

    def __str__(self):
        return "{}".format(self.name)

    def as_dict(self):
        return {
            'id': self.id,
            'name': self.name,
            'password': self.password,
            'counter': self.counter
        }


class Profile(models.Model):
    user = models.OneToOneField(AuthUser, unique=True, on_delete=models.CASCADE)
    public_key = models.CharField(max_length=30)


@receiver(post_save, sender=AuthUser)
def create_user_profile(sender, instance, created, **kwargs):
    if created:
        Profile.objects.create(user=instance)


@receiver(post_save, sender=AuthUser)
def save_user_profile(sender, instance, **kwargs):
    instance.profile.save()
