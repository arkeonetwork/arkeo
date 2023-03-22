{{ template "views/providers_base_v_v1.sql" . }}
{{ template "views/providers_v.sql" . }}
---- create above / drop below ----
drop view providers_v;
drop view providers_base_v;
