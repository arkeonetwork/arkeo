create or replace view providers_v as
(select b.*, b.cur_height - b.birth_height as age from providers_base_v b);
