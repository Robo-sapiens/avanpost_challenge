import numpy as np
import torch
from PIL import Image
from torch import nn


class NN(nn.Module):
    def __init__(self, input_size):
        super(NN, self).__init__()
        self.fc1 = nn.Linear(input_size, 600)

    def forward(self, x):
        x = self.fc1(x)
        return x


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
