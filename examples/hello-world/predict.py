from cog import BasePredictor, Input


class Predictor(BasePredictor):
    def setup(self):
        self.prefix = "hello"

    def predict(self, name: str = Input(description="What is your name?")) -> str:
        return self.prefix + " " + name