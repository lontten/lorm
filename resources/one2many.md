```sql
select 1 as a,
       2 as b,
       (
           select '[' || array_to_string(array(select row_to_json(a.*)
                                               from (
                                                        select b.is_head, b.app_name
                                                        from t_branch b
                                                    ) as a), ',') || ']'
       ) as c
;

```

```sql
select 1 as a,
       2 as b,
       (
           select '[' || array_to_string(array(select row_to_json(a.*)
                                               from (
                                                        ## sql_string
                                                    ) as a), ',') || ']'
       ) as ##name
;

```