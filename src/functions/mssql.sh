mssql_server_list_settings()
{
    mssql-conf list
}

mssql_server_set_setting()
{
    setting="${1}"

    mssql-conf set ${setting} true
}

mssql_server_set_default_setting()
{
    setting="${1}"

    mssql-conf unset ${setting}
}