from django.db import models

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
