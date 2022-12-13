create table svc_summ (
  id         serial primary key,
  block_id   integer,
  account    varchar(20) not null,
  sdate      date,
  code       varchar(20) not null,
  summ       numeric,
  info       varchar(255),
  created_at timestamp default current_timestamp   
);

create table pay_summ (
  id         serial primary key,
  block_id   integer,
  account    varchar(20) not null,
  pdate      date,
  summ       numeric,
  info       varchar(255),
  created_at timestamp default current_timestamp   
);
 