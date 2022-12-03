import random

import numpy as np
import torch


def seed(x):
    random.seed(x)
    np.random.seed(x)
    torch.manual_seed(x)
    torch.cuda.manual_seed(x)
    # Commented for faster training
    # torch.backends.cudnn.deterministic = True
