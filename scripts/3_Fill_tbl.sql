INSERT INTO block_type_ref(id, code, name, table_name) 
  VALUES (1, 'S','Начисления', 'svc_summ');
INSERT INTO block_type_ref(id, code, name, table_name) 
  VALUES (2, 'P','Оплаты', 'pay_summ');

-- S  
INSERT INTO block_type_dtl(type_id, field_name, field_alias, place, formatdate) VALUES (1, 'account', 'ЛС', 1, null);  
INSERT INTO block_type_dtl(type_id, field_name, field_alias, place, formatdate) VALUES (1, 'code', 'Код услуги', 2, null);  
INSERT INTO block_type_dtl(type_id, field_name, field_alias, place, formatdate) VALUES (1, 'summ', 'Сумма', 3, null); 
INSERT INTO block_type_dtl(type_id, field_name, field_alias, place, formatdate) VALUES (1, 'sdate', 'Дата', 4, 'DD.MM.YYYY'); 

-- P
INSERT INTO block_type_dtl(type_id, field_name, field_alias, place, formatdate) VALUES (2, 'account', 'ЛС', 1, null);  
INSERT INTO block_type_dtl(type_id, field_name, field_alias, place, formatdate) VALUES (2, 'pdate', 'Дата', 2, 'DD.MM.YYYY'); 
INSERT INTO block_type_dtl(type_id, field_name, field_alias, place, formatdate) VALUES (2, 'summ', 'Сумма', 3, null);
INSERT INTO block_type_dtl(type_id, field_name, field_alias, place, formatdate) VALUES (2, 'info', 'Инфо', 4, null);   

--
-- commit;


