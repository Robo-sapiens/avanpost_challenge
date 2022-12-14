import argparse
import os
from sys import stdin
from time import time
from typing import Any

from classifier import FingerprintClassifier


def _load_model() -> Any:
    path = os.environ.get("MODEL_PATH")
    if path is not None:
        return FingerprintClassifier(path)
    else:
        raise Exception("Environment variable 'MODEL_PATH' not found")


def cli_loop(model: Any, image_file: Any) -> None:
    pics_dir = os.environ.get("PICS")
    try:
        start = time()
        y = model.load_predict(os.path.join(pics_dir,image_file))
        print(y)
        end = time()
        print(f"time: {end - start}")
    except Exception as e:
        print(e)
#    for line in stdin:
#        line = line.strip()
#        if line == "exit":
#            print("Exiting")
#            return
#        try:
#            start = time()
#            y = model.load_predict(os.path.join(pics_dir,line))
#            print(y)
#            end = time()
#            print(f"time: {end - start}")
#        except Exception as e:
#            print(e)


def fifo_loop(model: Any, fifo: str) -> None:
    if not os.path.exists(fifo):
        raise Exception(f"Fifo with name {fifo} does not exist")

    with open(args.fifo, "r+") as fifo:
        path = fifo.read()

        start = time()
        y = model.predict(path)
        end = time()

        print(end - start)

        fifo.write(y)
        fifo.flush()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument("--fifo")
    group.add_argument("--cli", action="store_true")
    group.add_argument("--image_file", help="store filename which should be checked")
    args = parser.parse_args()
    try:
        model = _load_model()
        if args.fifo is not None:
            fifo_loop(model, args.fifo)
        elif args.cli is not None or args.image_file is not None:
            cli_loop(model, args.image_file)
    except Exception as e:
        print(e)
