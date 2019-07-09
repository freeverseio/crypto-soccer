from django.contrib.auth.tokens import PasswordResetTokenGenerator


class ValidationTokenGenerator(PasswordResetTokenGenerator):
    def _make_hash_value(self, user, timestamp):
        return (
                str(user.id)
                + str(timestamp)
                + str(user.is_active)
        )


account_activation_token = ValidationTokenGenerator()
