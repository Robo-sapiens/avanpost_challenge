import numpy as np
from PIL import Image


def parse_socofing(files, shape=(128, 128)):
    X = np.zeros(shape=(len(files), shape[0] * shape[1]))
    y = np.zeros(shape=(len(files)))
    for i, file in enumerate(files):
        image = Image.open(file).convert("L").resize(shape)
        X[i, :] = np.asarray(image).flatten()
        y[i] = int(file.name[0 : file.name.find("_")])
    return X, y
