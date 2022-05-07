INSERT INTO users (email, password, name, surname, "type")
VALUES ('test@test.com', 'vKrtXaX67RHf5zwpEdA1HtT8FxrhRRu6krsIUaZOTzo=', 'TESTA', 'TESTS', 1);
-- pass = 123
INSERT INTO users (email, password, name, surname, "type", created_by)
VALUES ('test2@test.com', 'pTzBptE5nGU6zFQi7uLYXx1W2iuJiQm2OjtI1xVFpzI=', 'AD', 'SOYAD', 0, 1);
-- pass = 123

INSERT INTO buildings (name, city_id, district_id, address)
VALUES ('Test Apartman 1', 43, 663, 'A Mh. C Blv. No:64');

INSERT INTO buildings (name, city_id, district_id, address)
VALUES ('Test Apartman 2', 43, 664, 'B Mh. D Blv. No:5');

INSERT INTO blocks (building_id, letter, d_number) VALUES (1, 'A', '64');
INSERT INTO blocks (building_id, letter, d_number) VALUES (1, 'B', '65');

INSERT INTO blocks (building_id, letter, d_number) VALUES (2, 'A', '5');

INSERT INTO user_building_membership (user_id, building_id, rank) VALUES (1, 1, 1);
INSERT INTO user_building_membership (user_id, building_id, rank) VALUES (1, 2, 1);

INSERT INTO flats (building_id, block_id, owner_id, tenant_id, number) VALUES (2, 3, 1, NULL, '1');
INSERT INTO flats (building_id, block_id, owner_id, tenant_id, number) VALUES (2, 3, NULL, NULL, '2');
INSERT INTO flats (building_id, block_id, owner_id, tenant_id, number) VALUES (2, 3, NULL, NULL, '3');
INSERT INTO flats (building_id, block_id, owner_id, tenant_id, number) VALUES (2, 3, NULL, NULL, '4');

INSERT INTO revenues (building_id, flat_id, rid, total, time, paid_type, paid_time, payer_full_name, payer_email, payer_phone, paid_status)
VALUES (2, 1, 1, 30.0, CURRENT_TIMESTAMP, 0, CURRENT_TIMESTAMP, 'TESTA TESTS', 'test@test.com', '+905555555555', true);
INSERT INTO revenues (building_id, flat_id, rid, total, time, paid_status)
VALUES (2, 1, 2, 30.5, CURRENT_TIMESTAMP, false);

