from PIL import Image
from torch.utils.data import Dataset
from torchvision import transforms


class SOCOFingDataset(Dataset):
    def __init__(self, files, size, mode):
        self.files = sorted(files)
        self.mode = mode
        self.transforms = transforms.Compose(
            [
                transforms.ToTensor(),
                transforms.Resize(size),
            ]
        )
        self.labels = [int(file.name[: file.name.find("_")]) - 1 for file in self.files]

    def __len__(self):
        return len(self.files)

    def __getitem__(self, index):
        x = self._load_sample(self.files[index])
        if self.transforms:
            x = self.transforms(x)
        y = self.labels[index]
        return x, y

    def _load_sample(self, file):
        image = Image.open(file)
        image = image.convert("L")
        return image
