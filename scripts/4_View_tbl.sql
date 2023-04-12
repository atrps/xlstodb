-- svc
select * from svc_summ order by id desc;
-- pay
select * from pay_summ order by id desc;

-- block
select * from block_type_ref order by id desc;
select * from block_type_dtl order by type_id desc, id asc;
select * from block_load order by id desc;


