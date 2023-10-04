from cog import BasePredictor


class Predictor(BasePredictor):
    def setup(self):
        self.prefix = "hello"

    def predict(self) -> str:
        return self.prefix + " world"