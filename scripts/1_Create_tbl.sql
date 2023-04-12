CREATE TABLE block_type_ref (
  id          serial PRIMARY KEY,
  code        varchar(10) CONSTRAINT uq_block_type_ref_code UNIQUE,
  name        varchar(100),
  table_name  varchar(60),
  info        varchar(255),
  created_at  timestamp default current_timestamp   
);

CREATE TABLE block_type_dtl (
  id          serial PRIMARY KEY,
  type_id     integer,
  field_name  varchar(60),
  field_alias varchar(100),
  place       integer,
  formatdate  varchar(20),
  FOREIGN KEY (type_id) REFERENCES block_type_ref(id) ON DELETE CASCADE
);

CREATE TABLE block_load (
  id          serial PRIMARY KEY,
  type_id     integer,
  info        varchar(255),
  created_at  timestamp default current_timestamp,
  status      varchar(1) default 'C',
  count_rec   integer,  
  completed_at timestamp,
  FOREIGN KEY (type_id) REFERENCES block_type_ref(id) ON DELETE RESTRICT  
);

CREATE TABLE svc_summ (
  id         serial PRIMARY KEY,
  block_id   integer,
  account    varchar(20) not null,
  sdate      date,
  code       varchar(20) not null,
  summ       numeric,
  info       varchar(255),
  created_at timestamp default current_timestamp
);
ALTER TABLE svc_summ
    ADD CONSTRAINT fk_svc_summ FOREIGN KEY (block_id) REFERENCES block_load(id) ON DELETE RESTRICT;

CREATE TABLE pay_summ (
  id         serial PRIMARY KEY,
  block_id   integer,
  account    varchar(20) not null,
  pdate      date,
  summ       numeric,
  info       varchar(255),
  created_at timestamp default current_timestamp
);
ALTER TABLE pay_summ
    ADD CONSTRAINT fk_pay_summ FOREIGN KEY (block_id) REFERENCES block_load(id) ON DELETE RESTRICT;