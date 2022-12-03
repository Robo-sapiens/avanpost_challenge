from PIL import Image
import torch
import numpy as np


class FingerprintClassifier:
    def __init__(self, path: str) -> None:
        self.model = torch.load(path)
        self.model.eval()
        self.device = torch.device("cuda:0")

    def predict(self, path: str) -> int:
        X = parse_fingerprints(path, shape=(64, 64))
        X = X.view(-1, 64 * 64).to(self.device)
        return self.model(X.view())


def parse_fingerprints(files, shape=(128, 128)):
    X = np.zeros(shape=(len(files), shape[0] * shape[1]))
    for i, file in enumerate(files):
        image = Image.open(file).convert("L").resize(shape)
        X[i, :] = np.asarray(image).flatten()
    return X
