import argparse
import os
from time import time

import torch

from classifier.socofing import parse_fingerprints


class FingerprintClassifier:
    def __init__(self, path: str) -> None:
        self.model = torch.load(path)
        self.model.eval()
        self.device = torch.device("cuda:0")

    def predict(self, path: str) -> int:
        X = parse_fingerprints(path, shape=(64, 64))
        X = X.view(-1, 64 * 64).to(self.device)
        return self.model(X.view())


def main(fifo: str) -> None:
    if not os.path.exists(fifo):
        raise Exception(f"Fifo with name {fifo} does not exist")

    model = FingerprintClassifier(os.path.join("models", "nn_sm_acc_99"))

    with open(args.fifo, "rw") as fifo:
        path = fifo.read()

        start = time()
        klass = model.predict(path)
        end = time()

        print(end - start)

        fifo.write(klass)
        fifo.flush()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--fifo", required=True)
    args = parser.parse_args()
    main(args.fifo)
