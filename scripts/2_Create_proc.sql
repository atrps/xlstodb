create or replace function block_create(p_type_code varchar,p_info varchar)
returns integer as $$
declare
  v_block_id integer;
  v_type_id  integer;
begin
  select id into v_type_id from block_type_ref where code = p_type_code;
  insert into block_load(type_id, info) values(v_type_id, p_info) returning id into v_block_id;
  return v_block_id;
exception
  when NO_DATA_FOUND then	
	return null;
  when OTHERS then
    return null;
end;
$$ language plpgsql;


create or replace procedure block_complete(p_block_id integer, p_cnt integer)
as $$
begin
  update block_load 
    set status = 'L', completed_at = current_timestamp, count_rec = p_cnt
    where id = p_block_id;
end;
$$ language plpgsql;

create or replace procedure block_change_status(p_block_id integer, p_status varchar default null , p_cnt integer default null)
as $$
begin
  update block_load b
    set status = coalesce(p_status, b.status), count_rec = coalesce(p_cnt, b.count_rec)
    where b.id = p_block_id;
end;
$$ language plpgsql;

create or replace function block_current_status(p_block_id integer)
returns varchar(100)[] as $$
declare
  v_status  varchar(1);
  v_strdate varchar; 
  v_cnt     varchar;
begin
  select status, count_rec::varchar as cnt, to_char(coalesce(completed_at, current_timestamp), 'DD.MM.YYYY HH24:MI:SS') as strdate  
    into v_status, v_cnt, v_strdate
    from block_load where id = p_block_id;
  return array[v_status, v_cnt, v_strdate];
exception
  when NO_DATA_FOUND then	
	return null;
end;
$$ language plpgsql;

create or replace function block_status_name(p_status_code varchar)
returns varchar as $$
declare
  v_name  varchar;
begin
 case p_status_code
   when 'C' then v_name := 'Создан';
   when 'L' then v_name := 'Загружен';
   when 'P' then v_name := 'В обработке'; 
   when 'E' then v_name := 'Ошибки при загрузке';    
 end case;   
 return v_name;
end;
$$ language plpgsql;

create or replace function block_type_dtl_fields(p_type_code varchar)
returns TABLE(dtl_fields json) as $$
begin
  return QUERY
    SELECT json_agg(json_build_object(
            'name', d.field_name,
            'place', d.place,
            'formatdate', d.formatdate
      )) dtl_fields
    FROM block_type_ref t
      INNER JOIN block_type_dtl d ON d.type_id= t.id 
    WHERE t.code = p_type_code; 
end;
$$ language plpgsql;







