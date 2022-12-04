import torch
from PIL import Image
from torch import nn
from torchvision import transforms


class NN(nn.Module):
    def __init__(self, input_size):
        super(NN, self).__init__()
        self.fc1 = nn.Linear(input_size, 600)

    def forward(self, x):
        x = self.fc1(x)
        return x


class FingerprintClassifier:
    def __init__(self, path: str) -> None:
        self.model = NN(64 * 64)
        self.model.load_state_dict(torch.load(path))
        self.model.eval()
        self.device = torch.device("cuda:0")
        self.model.to(self.device)

    def predict(self, X) -> int:
        y = self.model(X)
        y = torch.argmax(y, 1).cpu().numpy()[0]
        return y + 1

    def load_predict(self, path: str):
        X = _load_transform_image(path)
        X = X.view(-1, 64 * 64).to(self.device)
        return self.predict(X)


def _load_transform_image(path: str):
    image = Image.open(path)
    image = image.convert("L")
    t = transforms.Compose(
        [
            transforms.ToTensor(),
            transforms.Resize((64, 64)),
        ]
    )
    return t(image)
