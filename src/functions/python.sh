#!/bin/sh


venv()
{
    CWD=$(pwd)
    cd "${BASE_DIR}"
    
    python3 -m venv venv
    . venv/bin/activate
    pip3 install -r requirements.txt
    
    cd "${CWD}"
}

collectstatic()
{
    CWD=$(pwd)
    cd "${BASE_DIR}"

    python3 manage.py collectstatic --noinput

    cd "${CWD}"
}

migrate()
{
    CWD=$(pwd)
    cd "${BASE_DIR}"

    python3 manage.py makemigrations
    python3 manage.py migrate

    cd "${CWD}"
}

test()
{
    CWD=$(pwd)
    cd "${BASE_DIR}"
    
    python3 -m unitest -v

    cd "${CWD}"
}

runserver()
{
    CWD=$(pwd)
    cd "${BASE_DIR}"
    
    python3 manage.py collectstatic
    python3 manage.py runserver

    cd "${CWD}"
}
