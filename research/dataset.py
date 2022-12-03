from torch.utils.data import Dataset
from torchvision import transforms
from PIL import Image
from sklearn.preprocessing import LabelEncoder


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
        self.label_encoder = LabelEncoder()

        # if self.mode != "test":
        self.labels = [file.name[0 : file.name.find("_")] for file in self.files]
        self.label_encoder.fit(self.labels)

    def __len__(self):
        return len(self.files)

    def __getitem__(self, index):
        x = self._load_sample(self.files[index])
        if self.transforms:
            x = self.transforms(x)
        # if self.mode == "test":
        #     return x
        # else:
        y = self.labels[index]
        y = self.label_encoder.transform([y])
        y = y.item()
        return x, y

    def _load_sample(self, file):
        image = Image.open(file)
        image = image.convert("L")
        return image
