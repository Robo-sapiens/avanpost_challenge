import argparse
import os
from classifier import FingerprintClassifier
from time import time


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
