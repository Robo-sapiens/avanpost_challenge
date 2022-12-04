.PHONY: install format

install:
	pip install -r requirements.txt

format:
	python3 -m black .
	python3 -m isort --skip env .

curl:
	curl -X POST -H "Content-Type: application/json" -d @input.json http://127.0.0.1:10000/data
