#!/bin/sh


tf_init()
{
    if [ -n "$CI" ];
    then
        echo "Operation done in the backend by gitlab"
    else
        CWD=$(pwd)
        cd "${TF_ROOT}"

        terraform init \
            -backend-config="address=https://gitlab.com/api/v4/projects/${CI_PROJECT_ID}/terraform/state/${ENVIRONMENT}" \
            -backend-config="lock_address=https://gitlab.com/api/v4/projects/${CI_PROJECT_ID}/terraform/state/${ENVIRONMENT}/lock" \
            -backend-config="unlock_address=https://gitlab.com/api/v4/projects/${CI_PROJECT_ID}/terraform/state/${ENVIRONMENT}/lock" \
            -backend-config="username=TERRAFORM_BACKEND" \
            -backend-config="password=${TERRAFORM_PASSWD}" \
            -backend-config="lock_method=POST" \
            -backend-config="unlock_method=DELETE" \
            -backend-config="retry_wait_min=5"

        cd "${CWD}"
    fi
}

tf_validate()
{
    CWD=$(pwd)
    cd "${TF_ROOT}"

    if [ -n "$CI" ];
    then
        gitlab-terraform validate
    else
        terraform validate 
    fi

    cd "${CWD}"
}

tf_plan()
{
    CWD=$(pwd)
    cd "${TF_ROOT}"

    if [ -n "$CI" ];
    then
        if test "${REBUILD}" != "no";
        then
            replace=$(
            for i in $(gitlab-terraform state list module.ec2s);
            do
                echo "-replace=${i}"
            done
            )
            gitlab-terraform plan $replace
            gitlab-terraform plan-json $replace
        else
            gitlab-terraform plan
            gitlab-terraform plan-json
        fi;
    else
        terraform plan
        terraform plan-json
    fi

    cd "${CWD}"
}

tf_apply()
{
    CWD=$(pwd)
    cd "${TF_ROOT}"

    if [ -n "$CI" ];
    then
        gitlab-terraform apply
    else
        terraform apply
    fi

    cd "${CWD}"
}

tf_rebuild_ec2()
{
    CWD=$(pwd)
    cd "${TF_ROOT}"

    if [ -n "$CI" ];
    then
        echo "Unsupported operation"
        exit 1
    else
        terraform plan \
            -replace=module.ec2s[\"$1\"].aws_instance.default # e.g. CICDServer
            -target=module.ec2s
            -out my-plan
        terraform apply my-plan
    fi

    cd "${CWD}"
}

tf_destroy()
{
    CWD=$(pwd)
    cd "${TF_ROOT}"

    if [ -n "$CI" ];
    then
        gitlab-terraform destroy
    else
        terraform destroy 
    fi

    cd "${CWD}"
}
