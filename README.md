# Avanpost Challenge. Team Robo sapiens

## Team

* Vadim Sofin
* Sergey Khil

## Usage: 

```bash
export PICS=../SOCOFing/Real/ 
export MODEL_PATH="../models/nn_sm_acc_99
go run . --file ../input.json | python3 ../classifier/main.py --cli
```

to use web 
```bash
#server
PICS=../SOCOFing/Real/ MODEL_PATH="../models/nn_sm_acc_99" go run . -web

# on client
make curl
```

## Workspace setup

### Python env

```bash
Create env:
$ python3 -m venv env

Activate env:
$ source env/bin/activate

Exit env:
$ deactivate

Install dependencies:
$ pip3 install -r requirements.txt
```
